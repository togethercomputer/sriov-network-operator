// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	sriovnetworkv1 "github.com/togethercomputer/sriov-network-operator/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeSriovNetworkNodePolicies implements SriovNetworkNodePolicyInterface
type FakeSriovNetworkNodePolicies struct {
	Fake *FakeSriovnetworkV1
	ns   string
}

var sriovnetworknodepoliciesResource = schema.GroupVersionResource{Group: "sriovnetwork.openshift.io", Version: "v1", Resource: "sriovnetworknodepolicies"}

var sriovnetworknodepoliciesKind = schema.GroupVersionKind{Group: "sriovnetwork.openshift.io", Version: "v1", Kind: "SriovNetworkNodePolicy"}

// Get takes name of the sriovNetworkNodePolicy, and returns the corresponding sriovNetworkNodePolicy object, and an error if there is any.
func (c *FakeSriovNetworkNodePolicies) Get(ctx context.Context, name string, options v1.GetOptions) (result *sriovnetworkv1.SriovNetworkNodePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(sriovnetworknodepoliciesResource, c.ns, name), &sriovnetworkv1.SriovNetworkNodePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*sriovnetworkv1.SriovNetworkNodePolicy), err
}

// List takes label and field selectors, and returns the list of SriovNetworkNodePolicies that match those selectors.
func (c *FakeSriovNetworkNodePolicies) List(ctx context.Context, opts v1.ListOptions) (result *sriovnetworkv1.SriovNetworkNodePolicyList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(sriovnetworknodepoliciesResource, sriovnetworknodepoliciesKind, c.ns, opts), &sriovnetworkv1.SriovNetworkNodePolicyList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &sriovnetworkv1.SriovNetworkNodePolicyList{ListMeta: obj.(*sriovnetworkv1.SriovNetworkNodePolicyList).ListMeta}
	for _, item := range obj.(*sriovnetworkv1.SriovNetworkNodePolicyList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested sriovNetworkNodePolicies.
func (c *FakeSriovNetworkNodePolicies) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(sriovnetworknodepoliciesResource, c.ns, opts))

}

// Create takes the representation of a sriovNetworkNodePolicy and creates it.  Returns the server's representation of the sriovNetworkNodePolicy, and an error, if there is any.
func (c *FakeSriovNetworkNodePolicies) Create(ctx context.Context, sriovNetworkNodePolicy *sriovnetworkv1.SriovNetworkNodePolicy, opts v1.CreateOptions) (result *sriovnetworkv1.SriovNetworkNodePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(sriovnetworknodepoliciesResource, c.ns, sriovNetworkNodePolicy), &sriovnetworkv1.SriovNetworkNodePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*sriovnetworkv1.SriovNetworkNodePolicy), err
}

// Update takes the representation of a sriovNetworkNodePolicy and updates it. Returns the server's representation of the sriovNetworkNodePolicy, and an error, if there is any.
func (c *FakeSriovNetworkNodePolicies) Update(ctx context.Context, sriovNetworkNodePolicy *sriovnetworkv1.SriovNetworkNodePolicy, opts v1.UpdateOptions) (result *sriovnetworkv1.SriovNetworkNodePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(sriovnetworknodepoliciesResource, c.ns, sriovNetworkNodePolicy), &sriovnetworkv1.SriovNetworkNodePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*sriovnetworkv1.SriovNetworkNodePolicy), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeSriovNetworkNodePolicies) UpdateStatus(ctx context.Context, sriovNetworkNodePolicy *sriovnetworkv1.SriovNetworkNodePolicy, opts v1.UpdateOptions) (*sriovnetworkv1.SriovNetworkNodePolicy, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(sriovnetworknodepoliciesResource, "status", c.ns, sriovNetworkNodePolicy), &sriovnetworkv1.SriovNetworkNodePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*sriovnetworkv1.SriovNetworkNodePolicy), err
}

// Delete takes name of the sriovNetworkNodePolicy and deletes it. Returns an error if one occurs.
func (c *FakeSriovNetworkNodePolicies) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(sriovnetworknodepoliciesResource, c.ns, name), &sriovnetworkv1.SriovNetworkNodePolicy{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeSriovNetworkNodePolicies) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(sriovnetworknodepoliciesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &sriovnetworkv1.SriovNetworkNodePolicyList{})
	return err
}

// Patch applies the patch and returns the patched sriovNetworkNodePolicy.
func (c *FakeSriovNetworkNodePolicies) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *sriovnetworkv1.SriovNetworkNodePolicy, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(sriovnetworknodepoliciesResource, c.ns, name, pt, data, subresources...), &sriovnetworkv1.SriovNetworkNodePolicy{})

	if obj == nil {
		return nil, err
	}
	return obj.(*sriovnetworkv1.SriovNetworkNodePolicy), err
}
