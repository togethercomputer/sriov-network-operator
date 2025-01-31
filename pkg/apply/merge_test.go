package apply

import (
	"bytes"
	"testing"

	. "github.com/onsi/gomega"
	uns "k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// TestReconcileNamespace makes sure that namespace
// annotations are merged, and everything else is overwritten
// Namespaces use the "generic" logic; deployments and services
// have custom logic
func TestMergeNamespace(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: v1
kind: Namespace
metadata:
  name: ns1
  labels:
    a: cur
    b: cur
  annotations:
    a: cur
    b: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: v1
kind: Namespace
metadata:
  name: ns1
  labels:
    a: upd
    c: upd
  annotations:
    a: upd
    c: upd`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"a": "upd",
		"b": "cur",
		"c": "upd",
	}))

	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"a": "upd",
		"b": "cur",
		"c": "upd",
	}))
}

func TestMergeDeployment(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    a: cur
    b: cur
  annotations:
    deployment.kubernetes.io/revision: cur
    a: cur
    b: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    a: upd
    c: upd
  annotations:
    deployment.kubernetes.io/revision: upd
    a: upd
    c: upd`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	// labels are not merged
	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"a": "upd",
		"b": "cur",
		"c": "upd",
	}))

	// annotations are merged
	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"a": "upd",
		"b": "cur",
		"c": "upd",

		"deployment.kubernetes.io/revision": "cur",
	}))
}

func TestMergeOne(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    label-c: cur
  annotations:
    annotation-c: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"label-c": "cur",
	}))

	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"annotation-c": "cur",
	}))
}

func TestMergeNilCur(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    a: upd
    c: upd
  annotations:
    a: upd
    c: upd`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"a": "upd",
		"c": "upd",
	}))

	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"a": "upd",
		"c": "upd",
	}))
}

func TestMergeNilMeta(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(BeEmpty())
}

func TestMergeNilUpd(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1
  labels:
    a: cur
    b: cur
  annotations:
    a: cur
    b: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: d1`)

	// this mutates updated
	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(upd.GetLabels()).To(Equal(map[string]string{
		"a": "cur",
		"b": "cur",
	}))

	g.Expect(upd.GetAnnotations()).To(Equal(map[string]string{
		"a": "cur",
		"b": "cur",
	}))
}

func TestMergeService(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: v1
kind: Service
metadata:
  name: d1
spec:
  clusterIP: cur`)

	upd := UnstructuredFromYaml(t, `
apiVersion: v1
kind: Service
metadata:
  name: d1
spec:
  clusterIP: upd`)

	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	ip, _, err := uns.NestedString(upd.Object, "spec", "clusterIP")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ip).To(Equal("cur"))
}

func TestMergeServiceAccount(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: d1
  annotations:
    a: cur
secrets:
- foo`)

	upd := UnstructuredFromYaml(t, `
apiVersion: v1
kind: ServiceAccount
metadata:
  name: d1
  annotations:
    b: upd`)

	err := IsObjectSupported(cur)
	g.Expect(err).To(MatchError(ContainSubstring("cannot create ServiceAccount with secrets")))

	err = MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	s, ok, err := uns.NestedSlice(upd.Object, "secrets")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ok).To(BeTrue())
	g.Expect(s).To(ConsistOf("foo"))
}

func TestMergeWebHookCABundle(t *testing.T) {
	g := NewGomegaWithT(t)

	cur := UnstructuredFromYaml(t, `
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: webhook-1
  annotations:
    service.beta.openshift.io/inject-cabundle: "true"
webhooks:
  - name: webhook-name-1
    sideEffects: None
    admissionReviewVersions: ["v1", "v1beta1"]
    failurePolicy: Fail
    clientConfig:
      service:
        name: service1
        namespace: ns1
        path: "/path"
      caBundle: BASE64CABUNDLE
    rules:
      - operations: [ "CREATE", "UPDATE", "DELETE" ]
        apiGroups: ["sriovnetwork.openshift.io"]
        apiVersions: ["v1"]`)

	upd := UnstructuredFromYaml(t, `
apiVersion: admissionregistration.k8s.io/v1
kind: ValidatingWebhookConfiguration
metadata:
  name: webhook-1
  annotations:
    service.beta.openshift.io/inject-cabundle: "true"
webhooks:
  - name: webhook-name-1
    sideEffects: None
    admissionReviewVersions: ["v1", "v1beta1"]
    failurePolicy: Fail
    clientConfig:
      service:
        name: service1
        namespace: ns1
        path: "/path"
    rules:
      - operations: [ "CREATE", "UPDATE", "DELETE" ]
        apiGroups: ["sriovnetwork.openshift.io"]
        apiVersions: ["v1"]`)

	err := MergeObjectForUpdate(cur, upd)
	g.Expect(err).NotTo(HaveOccurred())

	webhooks, ok, err := uns.NestedSlice(upd.Object, "webhooks")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ok).To(BeTrue())
	g.Expect(len(webhooks)).To(Equal(1))

	webhook0, ok := webhooks[0].(map[string]interface{})
	g.Expect(ok).To(BeTrue())
	caBundle, ok, err := uns.NestedString(webhook0, "clientConfig", "caBundle")
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(ok).To(BeTrue())
	g.Expect(caBundle).To(Equal("BASE64CABUNDLE"))
}

// UnstructuredFromYaml creates an unstructured object from a raw yaml string
func UnstructuredFromYaml(t *testing.T, obj string) *uns.Unstructured {
	t.Helper()
	buf := bytes.NewBufferString(obj)
	decoder := yaml.NewYAMLOrJSONDecoder(buf, 4096)

	u := uns.Unstructured{}
	err := decoder.Decode(&u)
	if err != nil {
		t.Fatalf("failed to parse test yaml: %v", err)
	}

	return &u
}
