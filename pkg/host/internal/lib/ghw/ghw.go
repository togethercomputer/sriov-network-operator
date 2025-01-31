package ghw

import (
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/cpu"
)

func New() GHWLib {
	return &libWrapper{}
}

//go:generate ../../../../../bin/mockgen -destination mock/mock_ghw.go -source ghw.go
type GHWLib interface {
	// PCI returns a pointer to an Info that provide methods to access info about devices
	PCI() (Info, error)

	// CPU returns a pointer to an Info that provide methods to access info about devices
	CPU() (*cpu.Info, error)
}

// Info interface provide methods to access info about devices
type Info interface {
	// ListDevices returns a list of pointers to Device structs present on the
	// host system
	ListDevices() []*ghw.PCIDevice
}

type libWrapper struct{}

// PCI returns a pointer to an Info that provide methods to access info about devices
func (w *libWrapper) PCI() (Info, error) {
	return ghw.PCI()
}

func (w *libWrapper) CPU() (*cpu.Info, error) {
	return ghw.CPU()
}
