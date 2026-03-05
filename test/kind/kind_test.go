//go:build kind

package kind_test

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	dockertypes "github.com/docker/docker/api/types"
	dockerclient "github.com/docker/docker/client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/kind/pkg/apis/config/v1alpha4"
	"sigs.k8s.io/kind/pkg/cluster"
	"sigs.k8s.io/kind/pkg/cluster/nodes"
	"sigs.k8s.io/kind/pkg/cluster/nodeutils"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"

	sriovnetworkv1 "github.com/k8snetworkplumbingwg/sriov-network-operator/api/v1"
)

const (
	clusterName = "sriov-kind-test"
	namespace   = "sriov-network-operator"

	operatorImage     = "sriov-network-operator:test"
	configDaemonImage = "sriov-network-operator-config-daemon:test"
	webhookImage      = "sriov-network-operator-webhook:test"

	pollInterval = 5 * time.Second
	pollTimeout  = 10 * time.Minute
)

var (
	provider  *cluster.Provider
	k8sClient client.Client
	repoRoot  string
)

var _ = BeforeSuite(func() {
	var err error
	repoRoot, err = filepath.Abs(filepath.Join("..", ".."))
	Expect(err).NotTo(HaveOccurred())

	By("creating Kind cluster")
	provider = cluster.NewProvider()

	kindConfig := &v1alpha4.Cluster{
		Nodes: []v1alpha4.Node{
			{
				Role: v1alpha4.ControlPlaneRole,
				Labels: map[string]string{
					"node-role.kubernetes.io/worker": "",
				},
			},
		},
	}

	err = provider.Create(
		clusterName,
		cluster.CreateWithV1Alpha4Config(kindConfig),
		cluster.CreateWithNodeImage("kindest/node:v1.28.15"),
		cluster.CreateWithWaitForReady(5*time.Minute),
	)
	Expect(err).NotTo(HaveOccurred())

	By("getting kubeconfig")
	kubeconfigStr, err := provider.KubeConfig(clusterName, false)
	Expect(err).NotTo(HaveOccurred())

	kubeconfigPath := filepath.Join(GinkgoT().TempDir(), "kubeconfig")
	err = os.WriteFile(kubeconfigPath, []byte(kubeconfigStr), 0600)
	Expect(err).NotTo(HaveOccurred())

	By("building controller-runtime client")
	restCfg, err := clientcmd.RESTConfigFromKubeConfig([]byte(kubeconfigStr))
	Expect(err).NotTo(HaveOccurred())

	err = sriovnetworkv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = apiextensionsv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = admissionregistrationv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err = client.New(restCfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())

	By("building Docker images")
	ctx := context.Background()
	dockerCli, err := dockerclient.NewClientWithOpts(
		dockerclient.FromEnv,
		dockerclient.WithAPIVersionNegotiation(),
	)
	Expect(err).NotTo(HaveOccurred())
	defer dockerCli.Close()

	buildImage(ctx, dockerCli, repoRoot, "Dockerfile", operatorImage)
	buildImage(ctx, dockerCli, repoRoot, "Dockerfile.sriov-network-config-daemon", configDaemonImage)
	buildImage(ctx, dockerCli, repoRoot, "Dockerfile.webhook", webhookImage)

	By("loading images into Kind cluster")
	nodes, err := provider.ListNodes(clusterName)
	Expect(err).NotTo(HaveOccurred())
	Expect(nodes).NotTo(BeEmpty())

	for _, img := range []string{operatorImage, configDaemonImage, webhookImage} {
		loadImageIntoKind(ctx, dockerCli, nodes, img)
	}

	By("deploying cert-manager")
	deployCertManager(ctx, kubeconfigPath)

	By("installing sriov-network-operator Helm chart")
	installHelmChart(kubeconfigPath)

	By("waiting for SriovNetworkNodeState to reach Succeeded")
	Eventually(func(g Gomega) {
		nodeStateList := &sriovnetworkv1.SriovNetworkNodeStateList{}
		g.Expect(k8sClient.List(ctx, nodeStateList, client.InNamespace(namespace))).To(Succeed())
		g.Expect(nodeStateList.Items).NotTo(BeEmpty(), "no SriovNetworkNodeState objects found")
		for _, ns := range nodeStateList.Items {
			g.Expect(ns.Status.SyncStatus).To(Equal("Succeeded"),
				"node %s has SyncStatus %q", ns.Name, ns.Status.SyncStatus)
		}
	}).WithTimeout(pollTimeout).WithPolling(pollInterval).Should(Succeed())
})

var _ = AfterSuite(func() {
	if provider != nil {
		By("deleting Kind cluster")
		err := provider.Delete(clusterName, "")
		Expect(err).NotTo(HaveOccurred())
	}
})

var _ = Describe("Kind E2E", func() {
	ctx := context.Background()

	It("should have the operator Deployment available", func() {
		deploy := &appsv1.Deployment{}
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name: "sriov-network-operator", Namespace: namespace,
		}, deploy)
		Expect(err).NotTo(HaveOccurred())
		Expect(deploymentAvailable(deploy)).To(BeTrue(),
			"operator deployment not available: %+v", deploy.Status.Conditions)
	})

	It("should have config daemon DaemonSet pods running", func() {
		ds := &appsv1.DaemonSet{}
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name: "sriov-network-config-daemon", Namespace: namespace,
		}, ds)
		Expect(err).NotTo(HaveOccurred())
		Expect(ds.Status.NumberReady).To(BeNumerically(">", 0),
			"no ready config daemon pods")
		Expect(ds.Status.NumberReady).To(Equal(ds.Status.DesiredNumberScheduled))
	})

	It("should have CRDs registered", func() {
		expectedCRDs := []string{
			"sriovnetworknodepolicies.sriovnetwork.openshift.io",
			"sriovnetworknodestates.sriovnetwork.openshift.io",
			"sriovnetworks.sriovnetwork.openshift.io",
			"sriovoperatorconfigs.sriovnetwork.openshift.io",
			"sriovibnetworks.sriovnetwork.openshift.io",
			"sriovnetworkpoolconfigs.sriovnetwork.openshift.io",
			"ovsnetworks.sriovnetwork.openshift.io",
		}
		for _, name := range expectedCRDs {
			crd := &apiextensionsv1.CustomResourceDefinition{}
			err := k8sClient.Get(ctx, types.NamespacedName{Name: name}, crd)
			Expect(err).NotTo(HaveOccurred(), "CRD %s not found", name)
		}
	})

	It("should have SriovNetworkNodeState with SyncStatus Succeeded", func() {
		nodeStateList := &sriovnetworkv1.SriovNetworkNodeStateList{}
		err := k8sClient.List(ctx, nodeStateList, client.InNamespace(namespace))
		Expect(err).NotTo(HaveOccurred())
		Expect(nodeStateList.Items).NotTo(BeEmpty())
		for _, ns := range nodeStateList.Items {
			Expect(ns.Status.SyncStatus).To(Equal("Succeeded"))
		}
	})

	It("should have SriovOperatorConfig/default", func() {
		config := &sriovnetworkv1.SriovOperatorConfig{}
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name: "default", Namespace: namespace,
		}, config)
		Expect(err).NotTo(HaveOccurred())
	})

	It("should have webhook MutatingWebhookConfiguration", func() {
		mwc := &admissionregistrationv1.MutatingWebhookConfiguration{}
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name: "sriov-network-operator-webhook-config",
		}, mwc)
		Expect(err).NotTo(HaveOccurred())
		Expect(mwc.Webhooks).NotTo(BeEmpty())
	})

	It("should deploy the network-resources-injector DaemonSet", func() {
		config := &sriovnetworkv1.SriovOperatorConfig{}
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name: "default", Namespace: namespace,
		}, config)
		Expect(err).NotTo(HaveOccurred())

		ds := &appsv1.DaemonSet{}
		if config.Spec.EnableInjector {
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name: "network-resources-injector", Namespace: namespace,
			}, ds)
			Expect(err).NotTo(HaveOccurred())
			Expect(ds.Status.DesiredNumberScheduled).To(Equal(ds.Status.NumberReady))
		}
	})

	It("should deploy the operator-webhook DaemonSet", func() {
		config := &sriovnetworkv1.SriovOperatorConfig{}
		err := k8sClient.Get(ctx, types.NamespacedName{
			Name: "default", Namespace: namespace,
		}, config)
		Expect(err).NotTo(HaveOccurred())

		ds := &appsv1.DaemonSet{}
		if config.Spec.EnableOperatorWebhook {
			err = k8sClient.Get(ctx, types.NamespacedName{
				Name: "operator-webhook", Namespace: namespace,
			}, ds)
			Expect(err).NotTo(HaveOccurred())
			Expect(ds.Status.DesiredNumberScheduled).To(Equal(ds.Status.NumberReady))
		}
	})
})

// buildImage builds a Docker image from the given Dockerfile and context.
func buildImage(ctx context.Context, cli *dockerclient.Client, contextDir, dockerfile, tag string) {
	GinkgoHelper()

	buildCtx, err := createTarContext(contextDir)
	Expect(err).NotTo(HaveOccurred())

	resp, err := cli.ImageBuild(ctx, buildCtx, dockertypes.ImageBuildOptions{
		Tags:       []string{tag},
		Dockerfile: dockerfile,
		Remove:     true,
	})
	Expect(err).NotTo(HaveOccurred())
	defer resp.Body.Close()

	// Must read the response body to completion for the build to finish.
	_, err = io.Copy(GinkgoWriter, resp.Body)
	Expect(err).NotTo(HaveOccurred())
}

// createTarContext creates a tar archive of the directory for Docker build context.
func createTarContext(dir string) (io.Reader, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory and other unnecessary paths
		relPath, _ := filepath.Rel(dir, path)
		if strings.HasPrefix(relPath, ".git") || strings.HasPrefix(relPath, "Library") {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			header.Linkname = link
		}

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.Mode().IsRegular() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(tw, f)
		return err
	})
	if err != nil {
		return nil, err
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	return &buf, nil
}

// loadImageIntoKind saves a Docker image and loads it into all Kind nodes.
func loadImageIntoKind(ctx context.Context, cli *dockerclient.Client, kindNodes []nodes.Node, image string) {
	GinkgoHelper()

	imgReader, err := cli.ImageSave(ctx, []string{image})
	Expect(err).NotTo(HaveOccurred())
	defer imgReader.Close()

	for _, n := range kindNodes {
		err := nodeutils.LoadImageArchive(n, imgReader)
		Expect(err).NotTo(HaveOccurred())
	}
}

// deployCertManager installs cert-manager and waits for it to be ready.
func deployCertManager(ctx context.Context, kubeconfigPath string) {
	GinkgoHelper()

	certManagerURL := "https://github.com/cert-manager/cert-manager/releases/download/v1.12.0/cert-manager.yaml"

	// cert-manager YAML is 50k+ lines with CRDs that require server-side apply,
	// so we shell out to kubectl for this one step
	runCommand("kubectl", "--kubeconfig", kubeconfigPath, "apply", "-f", certManagerURL)

	By("waiting for cert-manager webhook to be ready")
	Eventually(func(g Gomega) {
		deploy := &appsv1.Deployment{}
		g.Expect(k8sClient.Get(ctx, types.NamespacedName{
			Name: "cert-manager-webhook", Namespace: "cert-manager",
		}, deploy)).To(Succeed())
		g.Expect(deploymentAvailable(deploy)).To(BeTrue())
	}).WithTimeout(3 * time.Minute).WithPolling(5 * time.Second).Should(Succeed())

	// Give cert-manager webhook a moment to start serving
	time.Sleep(10 * time.Second)
}

// installHelmChart installs the sriov-network-operator Helm chart.
func installHelmChart(kubeconfigPath string) {
	GinkgoHelper()

	helmCfg := new(action.Configuration)
	flags := genericclioptions.NewConfigFlags(true)
	flags.KubeConfig = &kubeconfigPath
	ns := namespace
	flags.Namespace = &ns

	err := helmCfg.Init(flags, namespace, "secrets", func(format string, v ...interface{}) {
		fmt.Fprintf(GinkgoWriter, format+"\n", v...)
	})
	Expect(err).NotTo(HaveOccurred())

	install := action.NewInstall(helmCfg)
	install.ReleaseName = "sriov-network-operator"
	install.Namespace = namespace
	install.CreateNamespace = true
	install.Wait = true
	install.Timeout = 5 * time.Minute

	chartPath := filepath.Join(repoRoot, "deployment", "sriov-network-operator-chart")
	chrt, err := loader.Load(chartPath)
	Expect(err).NotTo(HaveOccurred())

	values := map[string]interface{}{
		"images": map[string]interface{}{
			"operator":          operatorImage,
			"sriovConfigDaemon": configDaemonImage,
			"webhook":           webhookImage,
		},
		"operator": map[string]interface{}{
			"admissionControllers": map[string]interface{}{
				"enabled": true,
				"certificates": map[string]interface{}{
					"certManager": map[string]interface{}{
						"enabled":            true,
						"generateSelfSigned": true,
					},
				},
			},
		},
		"sriovOperatorConfig": map[string]interface{}{
			"deploy": true,
		},
	}

	_, err = install.Run(chrt, values)
	Expect(err).NotTo(HaveOccurred())
}

// runCommand executes a command, failing the test on error.
func runCommand(name string, args ...string) {
	GinkgoHelper()

	cmd := exec.Command(name, args...)
	cmd.Stdout = GinkgoWriter
	cmd.Stderr = GinkgoWriter
	err := cmd.Run()
	if err != nil {
		Fail(fmt.Sprintf("command %q failed: %v", name, err))
	}
}

func deploymentAvailable(deploy *appsv1.Deployment) bool {
	for _, c := range deploy.Status.Conditions {
		if c.Type == appsv1.DeploymentAvailable {
			return c.Status == corev1.ConditionTrue
		}
	}
	return false
}
