// Code generated by MockGen. DO NOT EDIT.
// Source: ./reconciler.go

// Package mock_input is a generated GoMock package.
package mock_input

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	v1alpha2 "github.com/solo-io/service-mesh-hub/pkg/api/certificates.smh.solo.io/v1alpha2"
	reconcile "github.com/solo-io/skv2/pkg/reconcile"
	v1 "k8s.io/api/core/v1"
)

// MockmultiClusterReconciler is a mock of multiClusterReconciler interface.
type MockmultiClusterReconciler struct {
	ctrl     *gomock.Controller
	recorder *MockmultiClusterReconcilerMockRecorder
}

// MockmultiClusterReconcilerMockRecorder is the mock recorder for MockmultiClusterReconciler.
type MockmultiClusterReconcilerMockRecorder struct {
	mock *MockmultiClusterReconciler
}

// NewMockmultiClusterReconciler creates a new mock instance.
func NewMockmultiClusterReconciler(ctrl *gomock.Controller) *MockmultiClusterReconciler {
	mock := &MockmultiClusterReconciler{ctrl: ctrl}
	mock.recorder = &MockmultiClusterReconcilerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockmultiClusterReconciler) EXPECT() *MockmultiClusterReconcilerMockRecorder {
	return m.recorder
}

// ReconcileIssuedCertificate mocks base method.
func (m *MockmultiClusterReconciler) ReconcileIssuedCertificate(clusterName string, obj *v1alpha2.IssuedCertificate) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileIssuedCertificate", clusterName, obj)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReconcileIssuedCertificate indicates an expected call of ReconcileIssuedCertificate.
func (mr *MockmultiClusterReconcilerMockRecorder) ReconcileIssuedCertificate(clusterName, obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileIssuedCertificate", reflect.TypeOf((*MockmultiClusterReconciler)(nil).ReconcileIssuedCertificate), clusterName, obj)
}

// ReconcileCertificateRequest mocks base method.
func (m *MockmultiClusterReconciler) ReconcileCertificateRequest(clusterName string, obj *v1alpha2.CertificateRequest) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileCertificateRequest", clusterName, obj)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReconcileCertificateRequest indicates an expected call of ReconcileCertificateRequest.
func (mr *MockmultiClusterReconcilerMockRecorder) ReconcileCertificateRequest(clusterName, obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileCertificateRequest", reflect.TypeOf((*MockmultiClusterReconciler)(nil).ReconcileCertificateRequest), clusterName, obj)
}

// ReconcileSecret mocks base method.
func (m *MockmultiClusterReconciler) ReconcileSecret(clusterName string, obj *v1.Secret) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileSecret", clusterName, obj)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReconcileSecret indicates an expected call of ReconcileSecret.
func (mr *MockmultiClusterReconcilerMockRecorder) ReconcileSecret(clusterName, obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileSecret", reflect.TypeOf((*MockmultiClusterReconciler)(nil).ReconcileSecret), clusterName, obj)
}

// MocksingleClusterReconciler is a mock of singleClusterReconciler interface.
type MocksingleClusterReconciler struct {
	ctrl     *gomock.Controller
	recorder *MocksingleClusterReconcilerMockRecorder
}

// MocksingleClusterReconcilerMockRecorder is the mock recorder for MocksingleClusterReconciler.
type MocksingleClusterReconcilerMockRecorder struct {
	mock *MocksingleClusterReconciler
}

// NewMocksingleClusterReconciler creates a new mock instance.
func NewMocksingleClusterReconciler(ctrl *gomock.Controller) *MocksingleClusterReconciler {
	mock := &MocksingleClusterReconciler{ctrl: ctrl}
	mock.recorder = &MocksingleClusterReconcilerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MocksingleClusterReconciler) EXPECT() *MocksingleClusterReconcilerMockRecorder {
	return m.recorder
}

// ReconcileIssuedCertificate mocks base method.
func (m *MocksingleClusterReconciler) ReconcileIssuedCertificate(obj *v1alpha2.IssuedCertificate) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileIssuedCertificate", obj)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReconcileIssuedCertificate indicates an expected call of ReconcileIssuedCertificate.
func (mr *MocksingleClusterReconcilerMockRecorder) ReconcileIssuedCertificate(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileIssuedCertificate", reflect.TypeOf((*MocksingleClusterReconciler)(nil).ReconcileIssuedCertificate), obj)
}

// ReconcileCertificateRequest mocks base method.
func (m *MocksingleClusterReconciler) ReconcileCertificateRequest(obj *v1alpha2.CertificateRequest) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileCertificateRequest", obj)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReconcileCertificateRequest indicates an expected call of ReconcileCertificateRequest.
func (mr *MocksingleClusterReconcilerMockRecorder) ReconcileCertificateRequest(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileCertificateRequest", reflect.TypeOf((*MocksingleClusterReconciler)(nil).ReconcileCertificateRequest), obj)
}

// ReconcileSecret mocks base method.
func (m *MocksingleClusterReconciler) ReconcileSecret(obj *v1.Secret) (reconcile.Result, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReconcileSecret", obj)
	ret0, _ := ret[0].(reconcile.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ReconcileSecret indicates an expected call of ReconcileSecret.
func (mr *MocksingleClusterReconcilerMockRecorder) ReconcileSecret(obj interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReconcileSecret", reflect.TypeOf((*MocksingleClusterReconciler)(nil).ReconcileSecret), obj)
}
