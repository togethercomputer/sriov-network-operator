// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"

	v1 "github.com/togethercomputer/sriov-network-operator/api/v1"
)

// SriovNetworkLister helps list SriovNetworks.
// All objects returned here must be treated as read-only.
type SriovNetworkLister interface {
	// List lists all SriovNetworks in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.SriovNetwork, err error)
	// SriovNetworks returns an object that can list and get SriovNetworks.
	SriovNetworks(namespace string) SriovNetworkNamespaceLister
	SriovNetworkListerExpansion
}

// sriovNetworkLister implements the SriovNetworkLister interface.
type sriovNetworkLister struct {
	indexer cache.Indexer
}

// NewSriovNetworkLister returns a new SriovNetworkLister.
func NewSriovNetworkLister(indexer cache.Indexer) SriovNetworkLister {
	return &sriovNetworkLister{indexer: indexer}
}

// List lists all SriovNetworks in the indexer.
func (s *sriovNetworkLister) List(selector labels.Selector) (ret []*v1.SriovNetwork, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.SriovNetwork))
	})
	return ret, err
}

// SriovNetworks returns an object that can list and get SriovNetworks.
func (s *sriovNetworkLister) SriovNetworks(namespace string) SriovNetworkNamespaceLister {
	return sriovNetworkNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// SriovNetworkNamespaceLister helps list and get SriovNetworks.
// All objects returned here must be treated as read-only.
type SriovNetworkNamespaceLister interface {
	// List lists all SriovNetworks in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.SriovNetwork, err error)
	// Get retrieves the SriovNetwork from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.SriovNetwork, error)
	SriovNetworkNamespaceListerExpansion
}

// sriovNetworkNamespaceLister implements the SriovNetworkNamespaceLister
// interface.
type sriovNetworkNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all SriovNetworks in the indexer for a given namespace.
func (s sriovNetworkNamespaceLister) List(selector labels.Selector) (ret []*v1.SriovNetwork, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.SriovNetwork))
	})
	return ret, err
}

// Get retrieves the SriovNetwork from the indexer for a given namespace and name.
func (s sriovNetworkNamespaceLister) Get(name string) (*v1.SriovNetwork, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("sriovnetwork"), name)
	}
	return obj.(*v1.SriovNetwork), nil
}
