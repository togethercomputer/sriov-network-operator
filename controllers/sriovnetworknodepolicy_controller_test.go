package controllers

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client/fake"

	dptypes "github.com/k8snetworkplumbingwg/sriov-network-device-plugin/pkg/types"

	sriovnetworkv1 "github.com/togethercomputer/sriov-network-operator/api/v1"
	v1 "github.com/togethercomputer/sriov-network-operator/api/v1"
	"github.com/togethercomputer/sriov-network-operator/pkg/consts"
	"github.com/togethercomputer/sriov-network-operator/pkg/featuregate"
	"github.com/togethercomputer/sriov-network-operator/pkg/vars"
)

func mustMarshallSelector(t *testing.T, input *dptypes.NetDeviceSelectors) *json.RawMessage {
	out, err := json.Marshal(input)
	if err != nil {
		t.Error(err)
		t.FailNow()
		return nil
	}
	ret := json.RawMessage(out)
	return &ret
}

func TestRenderDevicePluginConfigData(t *testing.T) {
	table := []struct {
		tname       string
		policy      sriovnetworkv1.SriovNetworkNodePolicy
		expResource dptypes.ResourceConfList
	}{
		{
			tname: "testVirtioVdpaVirtio",
			policy: sriovnetworkv1.SriovNetworkNodePolicy{
				Spec: v1.SriovNetworkNodePolicySpec{
					ResourceName: "resourceName",
					DeviceType:   consts.DeviceTypeNetDevice,
					VdpaType:     consts.VdpaTypeVirtio,
				},
			},
			expResource: dptypes.ResourceConfList{
				ResourceList: []dptypes.ResourceConfig{
					{
						ResourceName: "resourceName",
						Selectors: mustMarshallSelector(t, &dptypes.NetDeviceSelectors{
							VdpaType: dptypes.VdpaType(consts.VdpaTypeVirtio),
						}),
					},
				},
			},
		}, {
			tname: "testVhostVdpaVirtio",
			policy: sriovnetworkv1.SriovNetworkNodePolicy{
				Spec: v1.SriovNetworkNodePolicySpec{
					ResourceName: "resourceName",
					DeviceType:   consts.DeviceTypeNetDevice,
					VdpaType:     consts.VdpaTypeVhost,
				},
			},
			expResource: dptypes.ResourceConfList{
				ResourceList: []dptypes.ResourceConfig{
					{
						ResourceName: "resourceName",
						Selectors: mustMarshallSelector(t, &dptypes.NetDeviceSelectors{
							VdpaType: dptypes.VdpaType(consts.VdpaTypeVhost),
						}),
					},
				},
			},
		},
		{
			tname: "testExcludeTopology",
			policy: sriovnetworkv1.SriovNetworkNodePolicy{
				Spec: v1.SriovNetworkNodePolicySpec{
					ResourceName:    "resourceName",
					ExcludeTopology: true,
				},
			},
			expResource: dptypes.ResourceConfList{
				ResourceList: []dptypes.ResourceConfig{
					{
						ResourceName:    "resourceName",
						Selectors:       mustMarshallSelector(t, &dptypes.NetDeviceSelectors{}),
						ExcludeTopology: true,
					},
				},
			},
		},
	}

	reconciler := SriovNetworkNodePolicyReconciler{
		FeatureGate: featuregate.New(),
	}

	node := corev1.Node{ObjectMeta: metav1.ObjectMeta{Name: "node1"}}
	nodeState := sriovnetworkv1.SriovNetworkNodeState{ObjectMeta: metav1.ObjectMeta{Name: node.Name, Namespace: vars.Namespace}}

	scheme := runtime.NewScheme()
	utilruntime.Must(sriovnetworkv1.AddToScheme(scheme))
	reconciler.Client = fake.NewClientBuilder().
		WithScheme(scheme).WithObjects(&nodeState).
		Build()

	for _, tc := range table {
		policyList := sriovnetworkv1.SriovNetworkNodePolicyList{Items: []sriovnetworkv1.SriovNetworkNodePolicy{tc.policy}}

		t.Run(tc.tname, func(t *testing.T) {
			resourceList, err := reconciler.renderDevicePluginConfigData(context.TODO(), &policyList, &node)
			if err != nil {
				t.Error(tc.tname, "renderDevicePluginConfigData has failed")
			}

			if !cmp.Equal(resourceList, tc.expResource) {
				t.Error(tc.tname, "ResourceConfList not as expected", cmp.Diff(resourceList, tc.expResource))
			}
		})
	}
}
