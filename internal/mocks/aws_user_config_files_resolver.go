// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/recode-sh/cli/internal/aws (interfaces: UserConfigFilesResolver)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	userconfig "github.com/recode-sh/aws-cloud-provider/userconfig"
)

// AWSUserConfigFilesResolver is a mock of UserConfigFilesResolver interface.
type AWSUserConfigFilesResolver struct {
	ctrl     *gomock.Controller
	recorder *AWSUserConfigFilesResolverMockRecorder
}

// AWSUserConfigFilesResolverMockRecorder is the mock recorder for AWSUserConfigFilesResolver.
type AWSUserConfigFilesResolverMockRecorder struct {
	mock *AWSUserConfigFilesResolver
}

// NewAWSUserConfigFilesResolver creates a new mock instance.
func NewAWSUserConfigFilesResolver(ctrl *gomock.Controller) *AWSUserConfigFilesResolver {
	mock := &AWSUserConfigFilesResolver{ctrl: ctrl}
	mock.recorder = &AWSUserConfigFilesResolverMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *AWSUserConfigFilesResolver) EXPECT() *AWSUserConfigFilesResolverMockRecorder {
	return m.recorder
}

// Resolve mocks base method.
func (m *AWSUserConfigFilesResolver) Resolve() (*userconfig.Config, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Resolve")
	ret0, _ := ret[0].(*userconfig.Config)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Resolve indicates an expected call of Resolve.
func (mr *AWSUserConfigFilesResolverMockRecorder) Resolve() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Resolve", reflect.TypeOf((*AWSUserConfigFilesResolver)(nil).Resolve))
}
