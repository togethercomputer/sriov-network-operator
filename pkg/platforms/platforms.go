package platforms

import (
	"sigs.k8s.io/controller-runtime/pkg/log"

	"github.com/togethercomputer/sriov-network-operator/pkg/host"
	"github.com/togethercomputer/sriov-network-operator/pkg/platforms/openshift"
	"github.com/togethercomputer/sriov-network-operator/pkg/platforms/openstack"
	"github.com/togethercomputer/sriov-network-operator/pkg/utils"
)

//go:generate ../../bin/mockgen -destination mock/mock_platforms.go -source platforms.go
type Interface interface {
	openshift.OpenshiftContextInterface
	openstack.OpenstackInterface
}

type platformHelper struct {
	openshift.OpenshiftContextInterface
	openstack.OpenstackInterface
}

func NewDefaultPlatformHelper() (Interface, error) {
	openshiftContext, err := openshift.New()
	if err != nil {
		return nil, err
	}
	utilsHelper := utils.New()
	hostManager, err := host.NewHostManager(utilsHelper)
	if err != nil {
		log.Log.Error(err, "failed to create host manager")
		return nil, err
	}
	openstackContext := openstack.New(hostManager)

	return &platformHelper{
		openshiftContext,
		openstackContext,
	}, nil
}
