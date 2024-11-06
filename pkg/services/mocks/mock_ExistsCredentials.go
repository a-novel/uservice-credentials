// Code generated by mockery v2.46.3. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	services "github.com/a-novel/uservice-credentials/pkg/services"
	mock "github.com/stretchr/testify/mock"
)

// MockExistsCredentials is an autogenerated mock type for the ExistsCredentials type
type MockExistsCredentials struct {
	mock.Mock
}

type MockExistsCredentials_Expecter struct {
	mock *mock.Mock
}

func (_m *MockExistsCredentials) EXPECT() *MockExistsCredentials_Expecter {
	return &MockExistsCredentials_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, data
func (_m *MockExistsCredentials) Exec(ctx context.Context, data *services.ExistsCredentialsRequest) (*services.ExistsCredentialsResponse, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 *services.ExistsCredentialsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *services.ExistsCredentialsRequest) (*services.ExistsCredentialsResponse, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *services.ExistsCredentialsRequest) *services.ExistsCredentialsResponse); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*services.ExistsCredentialsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *services.ExistsCredentialsRequest) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockExistsCredentials_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockExistsCredentials_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - data *services.ExistsCredentialsRequest
func (_e *MockExistsCredentials_Expecter) Exec(ctx interface{}, data interface{}) *MockExistsCredentials_Exec_Call {
	return &MockExistsCredentials_Exec_Call{Call: _e.mock.On("Exec", ctx, data)}
}

func (_c *MockExistsCredentials_Exec_Call) Run(run func(ctx context.Context, data *services.ExistsCredentialsRequest)) *MockExistsCredentials_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*services.ExistsCredentialsRequest))
	})
	return _c
}

func (_c *MockExistsCredentials_Exec_Call) Return(_a0 *services.ExistsCredentialsResponse, _a1 error) *MockExistsCredentials_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockExistsCredentials_Exec_Call) RunAndReturn(run func(context.Context, *services.ExistsCredentialsRequest) (*services.ExistsCredentialsResponse, error)) *MockExistsCredentials_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockExistsCredentials creates a new instance of MockExistsCredentials. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockExistsCredentials(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockExistsCredentials {
	mock := &MockExistsCredentials{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
