// Code generated by MockGen. DO NOT EDIT.
// Source: platforms.go

// Package mock_platforms is a generated GoMock package.
package mock_platforms

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v10 "github.com/openshift/machine-config-operator/pkg/apis/machineconfiguration.openshift.io/v1"
	v1 "github.com/togethercomputer/sriov-network-operator/api/v1"
	openshift "github.com/togethercomputer/sriov-network-operator/pkg/platforms/openshift"
	v11 "k8s.io/api/core/v1"
)

// MockInterface is a mock of Interface interface.
type MockInterface struct {
	ctrl     *gomock.Controller
	recorder *MockInterfaceMockRecorder
}

// MockInterfaceMockRecorder is the mock recorder for MockInterface.
type MockInterfaceMockRecorder struct {
	mock *MockInterface
}

// NewMockInterface creates a new mock instance.
func NewMockInterface(ctrl *gomock.Controller) *MockInterface {
	mock := &MockInterface{ctrl: ctrl}
	mock.recorder = &MockInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInterface) EXPECT() *MockInterfaceMockRecorder {
	return m.recorder
}

// ChangeMachineConfigPoolPause mocks base method.
func (m *MockInterface) ChangeMachineConfigPoolPause(arg0 context.Context, arg1 *v10.MachineConfigPool, arg2 bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ChangeMachineConfigPoolPause", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ChangeMachineConfigPoolPause indicates an expected call of ChangeMachineConfigPoolPause.
func (mr *MockInterfaceMockRecorder) ChangeMachineConfigPoolPause(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChangeMachineConfigPoolPause", reflect.TypeOf((*MockInterface)(nil).ChangeMachineConfigPoolPause), arg0, arg1, arg2)
}

// CreateOpenstackDevicesInfo mocks base method.
func (m *MockInterface) CreateOpenstackDevicesInfo() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateOpenstackDevicesInfo")
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateOpenstackDevicesInfo indicates an expected call of CreateOpenstackDevicesInfo.
func (mr *MockInterfaceMockRecorder) CreateOpenstackDevicesInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOpenstackDevicesInfo", reflect.TypeOf((*MockInterface)(nil).CreateOpenstackDevicesInfo))
}

// CreateOpenstackDevicesInfoFromNodeStatus mocks base method.
func (m *MockInterface) CreateOpenstackDevicesInfoFromNodeStatus(arg0 *v1.SriovNetworkNodeState) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CreateOpenstackDevicesInfoFromNodeStatus", arg0)
}

// CreateOpenstackDevicesInfoFromNodeStatus indicates an expected call of CreateOpenstackDevicesInfoFromNodeStatus.
func (mr *MockInterfaceMockRecorder) CreateOpenstackDevicesInfoFromNodeStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateOpenstackDevicesInfoFromNodeStatus", reflect.TypeOf((*MockInterface)(nil).CreateOpenstackDevicesInfoFromNodeStatus), arg0)
}

// DiscoverSriovDevicesVirtual mocks base method.
func (m *MockInterface) DiscoverSriovDevicesVirtual() ([]v1.InterfaceExt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DiscoverSriovDevicesVirtual")
	ret0, _ := ret[0].([]v1.InterfaceExt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DiscoverSriovDevicesVirtual indicates an expected call of DiscoverSriovDevicesVirtual.
func (mr *MockInterfaceMockRecorder) DiscoverSriovDevicesVirtual() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DiscoverSriovDevicesVirtual", reflect.TypeOf((*MockInterface)(nil).DiscoverSriovDevicesVirtual))
}

// GetFlavor mocks base method.
func (m *MockInterface) GetFlavor() openshift.OpenshiftFlavor {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFlavor")
	ret0, _ := ret[0].(openshift.OpenshiftFlavor)
	return ret0
}

// GetFlavor indicates an expected call of GetFlavor.
func (mr *MockInterfaceMockRecorder) GetFlavor() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFlavor", reflect.TypeOf((*MockInterface)(nil).GetFlavor))
}

// GetNodeMachinePoolName mocks base method.
func (m *MockInterface) GetNodeMachinePoolName(arg0 context.Context, arg1 *v11.Node) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNodeMachinePoolName", arg0, arg1)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNodeMachinePoolName indicates an expected call of GetNodeMachinePoolName.
func (mr *MockInterfaceMockRecorder) GetNodeMachinePoolName(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNodeMachinePoolName", reflect.TypeOf((*MockInterface)(nil).GetNodeMachinePoolName), arg0, arg1)
}

// IsHypershift mocks base method.
func (m *MockInterface) IsHypershift() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsHypershift")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsHypershift indicates an expected call of IsHypershift.
func (mr *MockInterfaceMockRecorder) IsHypershift() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsHypershift", reflect.TypeOf((*MockInterface)(nil).IsHypershift))
}

// IsOpenshiftCluster mocks base method.
func (m *MockInterface) IsOpenshiftCluster() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsOpenshiftCluster")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsOpenshiftCluster indicates an expected call of IsOpenshiftCluster.
func (mr *MockInterfaceMockRecorder) IsOpenshiftCluster() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsOpenshiftCluster", reflect.TypeOf((*MockInterface)(nil).IsOpenshiftCluster))
}

// OpenshiftAfterCompleteDrainNode mocks base method.
func (m *MockInterface) OpenshiftAfterCompleteDrainNode(arg0 context.Context, arg1 *v11.Node) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenshiftAfterCompleteDrainNode", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenshiftAfterCompleteDrainNode indicates an expected call of OpenshiftAfterCompleteDrainNode.
func (mr *MockInterfaceMockRecorder) OpenshiftAfterCompleteDrainNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenshiftAfterCompleteDrainNode", reflect.TypeOf((*MockInterface)(nil).OpenshiftAfterCompleteDrainNode), arg0, arg1)
}

// OpenshiftBeforeDrainNode mocks base method.
func (m *MockInterface) OpenshiftBeforeDrainNode(arg0 context.Context, arg1 *v11.Node) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "OpenshiftBeforeDrainNode", arg0, arg1)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// OpenshiftBeforeDrainNode indicates an expected call of OpenshiftBeforeDrainNode.
func (mr *MockInterfaceMockRecorder) OpenshiftBeforeDrainNode(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "OpenshiftBeforeDrainNode", reflect.TypeOf((*MockInterface)(nil).OpenshiftBeforeDrainNode), arg0, arg1)
}
