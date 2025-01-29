package host

import (
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/bridge"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/infiniband"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/kernel"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/lib/dputils"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/lib/ethtool"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/lib/ghw"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/lib/netlink"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/lib/sriovnet"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/network"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/service"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/sriov"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/udev"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/internal/vdpa"
	"github.com/togethercomputer/sriov-network-operator/pkg/host/types"
	"github.com/togethercomputer/sriov-network-operator/pkg/utils"
)

// Contains all the host manipulation functions
//
//go:generate ../../bin/mockgen -destination mock/mock_host.go -source manager.go
type HostManagerInterface interface {
	types.KernelInterface
	types.NetworkInterface
	types.ServiceInterface
	types.UdevInterface
	types.SriovInterface
	types.VdpaInterface
	types.InfinibandInterface
	types.BridgeInterface
}

type hostManager struct {
	utils.CmdInterface
	types.KernelInterface
	types.NetworkInterface
	types.ServiceInterface
	types.UdevInterface
	types.SriovInterface
	types.VdpaInterface
	types.InfinibandInterface
	types.BridgeInterface
}

func NewHostManager(utilsInterface utils.CmdInterface) (HostManagerInterface, error) {
	dpUtils := dputils.New()
	netlinkLib := netlink.New()
	ethtoolLib := ethtool.New()
	sriovnetLib := sriovnet.New()
	ghwLib := ghw.New()
	k := kernel.New(utilsInterface)
	n := network.New(utilsInterface, dpUtils, netlinkLib, ethtoolLib)
	sv := service.New(utilsInterface)
	u := udev.New(utilsInterface)
	v := vdpa.New(k, netlinkLib)
	ib, err := infiniband.New(netlinkLib, k, n)
	if err != nil {
		return nil, err
	}
	br := bridge.New()
	sr := sriov.New(utilsInterface, k, n, u, v, ib, netlinkLib, dpUtils, sriovnetLib, ghwLib, br)
	return &hostManager{
		utilsInterface,
		k,
		n,
		sv,
		u,
		sr,
		v,
		ib,
		br,
	}, nil
}
