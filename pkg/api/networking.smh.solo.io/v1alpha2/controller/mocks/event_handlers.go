// Code generated by MockGen. DO NOT EDIT.
// Source: ./event_handlers.go

// Package mock_controller is a generated GoMock package.
package mock_controller

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha2 "github.com/solo-io/service-mesh-hub/pkg/api/networking.smh.solo.io/v1alpha2"
	controller "github.com/solo-io/service-mesh-hub/pkg/api/networking.smh.solo.io/v1alpha2/controller"
	predicate "sigs.k8s.io/controller-runtime/pkg/predicate"
)

// MockTrafficPolicyEventHandler is a mock of TrafficPolicyEventHandler interface.
type MockTrafficPolicyEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockTrafficPolicyEventHandlerMockRecorder
}

// MockTrafficPolicyEventHandlerMockRecorder is the mock recorder for MockTrafficPolicyEventHandler.
type MockTrafficPolicyEventHandlerMockRecorder struct {
	mock *MockTrafficPolicyEventHandler
}

// NewMockTrafficPolicyEventHandler creates a new mock instance.
func NewMockTrafficPolicyEventHandler(ctrl *gomock.Controller) *MockTrafficPolicyEventHandler {
	mock := &MockTrafficPolicyEventHandler{ctrl: ctrl}
	mock.recorder = &MockTrafficPolicyEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrafficPolicyEventHandler) EXPECT() *MockTrafficPolicyEventHandlerMockRecorder {
	return m.recorder
}

// CreateTrafficPolicy mocks base method.
func (m *MockTrafficPolicyEventHandler) CreateTrafficPolicy(obj *v1alpha2.TrafficPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTrafficPolicy", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTrafficPolicy indicates an expected call of CreateTrafficPolicy.
func (mr *MockTrafficPolicyEventHandlerMockRecorder) CreateTrafficPolicy(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTrafficPolicy", reflect.TypeOf((*MockTrafficPolicyEventHandler)(nil).CreateTrafficPolicy), obj)
}

// UpdateTrafficPolicy mocks base method.
func (m *MockTrafficPolicyEventHandler) UpdateTrafficPolicy(old, new *v1alpha2.TrafficPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTrafficPolicy", old, new)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTrafficPolicy indicates an expected call of UpdateTrafficPolicy.
func (mr *MockTrafficPolicyEventHandlerMockRecorder) UpdateTrafficPolicy(old, new interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTrafficPolicy", reflect.TypeOf((*MockTrafficPolicyEventHandler)(nil).UpdateTrafficPolicy), old, new)
}

// DeleteTrafficPolicy mocks base method.
func (m *MockTrafficPolicyEventHandler) DeleteTrafficPolicy(obj *v1alpha2.TrafficPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteTrafficPolicy", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteTrafficPolicy indicates an expected call of DeleteTrafficPolicy.
func (mr *MockTrafficPolicyEventHandlerMockRecorder) DeleteTrafficPolicy(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteTrafficPolicy", reflect.TypeOf((*MockTrafficPolicyEventHandler)(nil).DeleteTrafficPolicy), obj)
}

// GenericTrafficPolicy mocks base method.
func (m *MockTrafficPolicyEventHandler) GenericTrafficPolicy(obj *v1alpha2.TrafficPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenericTrafficPolicy", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenericTrafficPolicy indicates an expected call of GenericTrafficPolicy.
func (mr *MockTrafficPolicyEventHandlerMockRecorder) GenericTrafficPolicy(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenericTrafficPolicy", reflect.TypeOf((*MockTrafficPolicyEventHandler)(nil).GenericTrafficPolicy), obj)
}

// MockTrafficPolicyEventWatcher is a mock of TrafficPolicyEventWatcher interface.
type MockTrafficPolicyEventWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockTrafficPolicyEventWatcherMockRecorder
}

// MockTrafficPolicyEventWatcherMockRecorder is the mock recorder for MockTrafficPolicyEventWatcher.
type MockTrafficPolicyEventWatcherMockRecorder struct {
	mock *MockTrafficPolicyEventWatcher
}

// NewMockTrafficPolicyEventWatcher creates a new mock instance.
func NewMockTrafficPolicyEventWatcher(ctrl *gomock.Controller) *MockTrafficPolicyEventWatcher {
	mock := &MockTrafficPolicyEventWatcher{ctrl: ctrl}
	mock.recorder = &MockTrafficPolicyEventWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTrafficPolicyEventWatcher) EXPECT() *MockTrafficPolicyEventWatcherMockRecorder {
	return m.recorder
}

// AddEventHandler mocks base method.
func (m *MockTrafficPolicyEventWatcher) AddEventHandler(ctx context.Context, h controller.TrafficPolicyEventHandler, predicates ...predicate.Predicate) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, h}
	for _, a := range predicates {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddEventHandler", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEventHandler indicates an expected call of AddEventHandler.
func (mr *MockTrafficPolicyEventWatcherMockRecorder) AddEventHandler(ctx, h interface{}, predicates ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, h}, predicates...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventHandler", reflect.TypeOf((*MockTrafficPolicyEventWatcher)(nil).AddEventHandler), varargs...)
}

// MockAccessPolicyEventHandler is a mock of AccessPolicyEventHandler interface.
type MockAccessPolicyEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockAccessPolicyEventHandlerMockRecorder
}

// MockAccessPolicyEventHandlerMockRecorder is the mock recorder for MockAccessPolicyEventHandler.
type MockAccessPolicyEventHandlerMockRecorder struct {
	mock *MockAccessPolicyEventHandler
}

// NewMockAccessPolicyEventHandler creates a new mock instance.
func NewMockAccessPolicyEventHandler(ctrl *gomock.Controller) *MockAccessPolicyEventHandler {
	mock := &MockAccessPolicyEventHandler{ctrl: ctrl}
	mock.recorder = &MockAccessPolicyEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccessPolicyEventHandler) EXPECT() *MockAccessPolicyEventHandlerMockRecorder {
	return m.recorder
}

// CreateAccessPolicy mocks base method.
func (m *MockAccessPolicyEventHandler) CreateAccessPolicy(obj *v1alpha2.AccessPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateAccessPolicy", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateAccessPolicy indicates an expected call of CreateAccessPolicy.
func (mr *MockAccessPolicyEventHandlerMockRecorder) CreateAccessPolicy(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateAccessPolicy", reflect.TypeOf((*MockAccessPolicyEventHandler)(nil).CreateAccessPolicy), obj)
}

// UpdateAccessPolicy mocks base method.
func (m *MockAccessPolicyEventHandler) UpdateAccessPolicy(old, new *v1alpha2.AccessPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateAccessPolicy", old, new)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateAccessPolicy indicates an expected call of UpdateAccessPolicy.
func (mr *MockAccessPolicyEventHandlerMockRecorder) UpdateAccessPolicy(old, new interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateAccessPolicy", reflect.TypeOf((*MockAccessPolicyEventHandler)(nil).UpdateAccessPolicy), old, new)
}

// DeleteAccessPolicy mocks base method.
func (m *MockAccessPolicyEventHandler) DeleteAccessPolicy(obj *v1alpha2.AccessPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteAccessPolicy", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAccessPolicy indicates an expected call of DeleteAccessPolicy.
func (mr *MockAccessPolicyEventHandlerMockRecorder) DeleteAccessPolicy(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAccessPolicy", reflect.TypeOf((*MockAccessPolicyEventHandler)(nil).DeleteAccessPolicy), obj)
}

// GenericAccessPolicy mocks base method.
func (m *MockAccessPolicyEventHandler) GenericAccessPolicy(obj *v1alpha2.AccessPolicy) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenericAccessPolicy", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenericAccessPolicy indicates an expected call of GenericAccessPolicy.
func (mr *MockAccessPolicyEventHandlerMockRecorder) GenericAccessPolicy(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenericAccessPolicy", reflect.TypeOf((*MockAccessPolicyEventHandler)(nil).GenericAccessPolicy), obj)
}

// MockAccessPolicyEventWatcher is a mock of AccessPolicyEventWatcher interface.
type MockAccessPolicyEventWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockAccessPolicyEventWatcherMockRecorder
}

// MockAccessPolicyEventWatcherMockRecorder is the mock recorder for MockAccessPolicyEventWatcher.
type MockAccessPolicyEventWatcherMockRecorder struct {
	mock *MockAccessPolicyEventWatcher
}

// NewMockAccessPolicyEventWatcher creates a new mock instance.
func NewMockAccessPolicyEventWatcher(ctrl *gomock.Controller) *MockAccessPolicyEventWatcher {
	mock := &MockAccessPolicyEventWatcher{ctrl: ctrl}
	mock.recorder = &MockAccessPolicyEventWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccessPolicyEventWatcher) EXPECT() *MockAccessPolicyEventWatcherMockRecorder {
	return m.recorder
}

// AddEventHandler mocks base method.
func (m *MockAccessPolicyEventWatcher) AddEventHandler(ctx context.Context, h controller.AccessPolicyEventHandler, predicates ...predicate.Predicate) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, h}
	for _, a := range predicates {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddEventHandler", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEventHandler indicates an expected call of AddEventHandler.
func (mr *MockAccessPolicyEventWatcherMockRecorder) AddEventHandler(ctx, h interface{}, predicates ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, h}, predicates...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventHandler", reflect.TypeOf((*MockAccessPolicyEventWatcher)(nil).AddEventHandler), varargs...)
}

// MockVirtualMeshEventHandler is a mock of VirtualMeshEventHandler interface.
type MockVirtualMeshEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockVirtualMeshEventHandlerMockRecorder
}

// MockVirtualMeshEventHandlerMockRecorder is the mock recorder for MockVirtualMeshEventHandler.
type MockVirtualMeshEventHandlerMockRecorder struct {
	mock *MockVirtualMeshEventHandler
}

// NewMockVirtualMeshEventHandler creates a new mock instance.
func NewMockVirtualMeshEventHandler(ctrl *gomock.Controller) *MockVirtualMeshEventHandler {
	mock := &MockVirtualMeshEventHandler{ctrl: ctrl}
	mock.recorder = &MockVirtualMeshEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVirtualMeshEventHandler) EXPECT() *MockVirtualMeshEventHandlerMockRecorder {
	return m.recorder
}

// CreateVirtualMesh mocks base method.
func (m *MockVirtualMeshEventHandler) CreateVirtualMesh(obj *v1alpha2.VirtualMesh) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateVirtualMesh", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateVirtualMesh indicates an expected call of CreateVirtualMesh.
func (mr *MockVirtualMeshEventHandlerMockRecorder) CreateVirtualMesh(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateVirtualMesh", reflect.TypeOf((*MockVirtualMeshEventHandler)(nil).CreateVirtualMesh), obj)
}

// UpdateVirtualMesh mocks base method.
func (m *MockVirtualMeshEventHandler) UpdateVirtualMesh(old, new *v1alpha2.VirtualMesh) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVirtualMesh", old, new)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateVirtualMesh indicates an expected call of UpdateVirtualMesh.
func (mr *MockVirtualMeshEventHandlerMockRecorder) UpdateVirtualMesh(old, new interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVirtualMesh", reflect.TypeOf((*MockVirtualMeshEventHandler)(nil).UpdateVirtualMesh), old, new)
}

// DeleteVirtualMesh mocks base method.
func (m *MockVirtualMeshEventHandler) DeleteVirtualMesh(obj *v1alpha2.VirtualMesh) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteVirtualMesh", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteVirtualMesh indicates an expected call of DeleteVirtualMesh.
func (mr *MockVirtualMeshEventHandlerMockRecorder) DeleteVirtualMesh(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteVirtualMesh", reflect.TypeOf((*MockVirtualMeshEventHandler)(nil).DeleteVirtualMesh), obj)
}

// GenericVirtualMesh mocks base method.
func (m *MockVirtualMeshEventHandler) GenericVirtualMesh(obj *v1alpha2.VirtualMesh) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenericVirtualMesh", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenericVirtualMesh indicates an expected call of GenericVirtualMesh.
func (mr *MockVirtualMeshEventHandlerMockRecorder) GenericVirtualMesh(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenericVirtualMesh", reflect.TypeOf((*MockVirtualMeshEventHandler)(nil).GenericVirtualMesh), obj)
}

// MockVirtualMeshEventWatcher is a mock of VirtualMeshEventWatcher interface.
type MockVirtualMeshEventWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockVirtualMeshEventWatcherMockRecorder
}

// MockVirtualMeshEventWatcherMockRecorder is the mock recorder for MockVirtualMeshEventWatcher.
type MockVirtualMeshEventWatcherMockRecorder struct {
	mock *MockVirtualMeshEventWatcher
}

// NewMockVirtualMeshEventWatcher creates a new mock instance.
func NewMockVirtualMeshEventWatcher(ctrl *gomock.Controller) *MockVirtualMeshEventWatcher {
	mock := &MockVirtualMeshEventWatcher{ctrl: ctrl}
	mock.recorder = &MockVirtualMeshEventWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVirtualMeshEventWatcher) EXPECT() *MockVirtualMeshEventWatcherMockRecorder {
	return m.recorder
}

// AddEventHandler mocks base method.
func (m *MockVirtualMeshEventWatcher) AddEventHandler(ctx context.Context, h controller.VirtualMeshEventHandler, predicates ...predicate.Predicate) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, h}
	for _, a := range predicates {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddEventHandler", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEventHandler indicates an expected call of AddEventHandler.
func (mr *MockVirtualMeshEventWatcherMockRecorder) AddEventHandler(ctx, h interface{}, predicates ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, h}, predicates...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventHandler", reflect.TypeOf((*MockVirtualMeshEventWatcher)(nil).AddEventHandler), varargs...)
}

// MockFailoverServiceEventHandler is a mock of FailoverServiceEventHandler interface.
type MockFailoverServiceEventHandler struct {
	ctrl     *gomock.Controller
	recorder *MockFailoverServiceEventHandlerMockRecorder
}

// MockFailoverServiceEventHandlerMockRecorder is the mock recorder for MockFailoverServiceEventHandler.
type MockFailoverServiceEventHandlerMockRecorder struct {
	mock *MockFailoverServiceEventHandler
}

// NewMockFailoverServiceEventHandler creates a new mock instance.
func NewMockFailoverServiceEventHandler(ctrl *gomock.Controller) *MockFailoverServiceEventHandler {
	mock := &MockFailoverServiceEventHandler{ctrl: ctrl}
	mock.recorder = &MockFailoverServiceEventHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFailoverServiceEventHandler) EXPECT() *MockFailoverServiceEventHandlerMockRecorder {
	return m.recorder
}

// CreateFailoverService mocks base method.
func (m *MockFailoverServiceEventHandler) CreateFailoverService(obj *v1alpha2.FailoverService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFailoverService", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFailoverService indicates an expected call of CreateFailoverService.
func (mr *MockFailoverServiceEventHandlerMockRecorder) CreateFailoverService(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFailoverService", reflect.TypeOf((*MockFailoverServiceEventHandler)(nil).CreateFailoverService), obj)
}

// UpdateFailoverService mocks base method.
func (m *MockFailoverServiceEventHandler) UpdateFailoverService(old, new *v1alpha2.FailoverService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFailoverService", old, new)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFailoverService indicates an expected call of UpdateFailoverService.
func (mr *MockFailoverServiceEventHandlerMockRecorder) UpdateFailoverService(old, new interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFailoverService", reflect.TypeOf((*MockFailoverServiceEventHandler)(nil).UpdateFailoverService), old, new)
}

// DeleteFailoverService mocks base method.
func (m *MockFailoverServiceEventHandler) DeleteFailoverService(obj *v1alpha2.FailoverService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFailoverService", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFailoverService indicates an expected call of DeleteFailoverService.
func (mr *MockFailoverServiceEventHandlerMockRecorder) DeleteFailoverService(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFailoverService", reflect.TypeOf((*MockFailoverServiceEventHandler)(nil).DeleteFailoverService), obj)
}

// GenericFailoverService mocks base method.
func (m *MockFailoverServiceEventHandler) GenericFailoverService(obj *v1alpha2.FailoverService) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenericFailoverService", obj)
	ret0, _ := ret[0].(error)
	return ret0
}

// GenericFailoverService indicates an expected call of GenericFailoverService.
func (mr *MockFailoverServiceEventHandlerMockRecorder) GenericFailoverService(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenericFailoverService", reflect.TypeOf((*MockFailoverServiceEventHandler)(nil).GenericFailoverService), obj)
}

// MockFailoverServiceEventWatcher is a mock of FailoverServiceEventWatcher interface.
type MockFailoverServiceEventWatcher struct {
	ctrl     *gomock.Controller
	recorder *MockFailoverServiceEventWatcherMockRecorder
}

// MockFailoverServiceEventWatcherMockRecorder is the mock recorder for MockFailoverServiceEventWatcher.
type MockFailoverServiceEventWatcherMockRecorder struct {
	mock *MockFailoverServiceEventWatcher
}

// NewMockFailoverServiceEventWatcher creates a new mock instance.
func NewMockFailoverServiceEventWatcher(ctrl *gomock.Controller) *MockFailoverServiceEventWatcher {
	mock := &MockFailoverServiceEventWatcher{ctrl: ctrl}
	mock.recorder = &MockFailoverServiceEventWatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFailoverServiceEventWatcher) EXPECT() *MockFailoverServiceEventWatcherMockRecorder {
	return m.recorder
}

// AddEventHandler mocks base method.
func (m *MockFailoverServiceEventWatcher) AddEventHandler(ctx context.Context, h controller.FailoverServiceEventHandler, predicates ...predicate.Predicate) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, h}
	for _, a := range predicates {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddEventHandler", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddEventHandler indicates an expected call of AddEventHandler.
func (mr *MockFailoverServiceEventWatcherMockRecorder) AddEventHandler(ctx, h interface{}, predicates ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, h}, predicates...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddEventHandler", reflect.TypeOf((*MockFailoverServiceEventWatcher)(nil).AddEventHandler), varargs...)
}
