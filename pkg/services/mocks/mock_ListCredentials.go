// Code generated by mockery v2.46.3. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	services "github.com/a-novel/uservice-credentials/pkg/services"
	mock "github.com/stretchr/testify/mock"
)

// MockListCredentials is an autogenerated mock type for the ListCredentials type
type MockListCredentials struct {
	mock.Mock
}

type MockListCredentials_Expecter struct {
	mock *mock.Mock
}

func (_m *MockListCredentials) EXPECT() *MockListCredentials_Expecter {
	return &MockListCredentials_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, data
func (_m *MockListCredentials) Exec(ctx context.Context, data *services.ListCredentialsRequest) (*services.ListCredentialsResponse, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 *services.ListCredentialsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *services.ListCredentialsRequest) (*services.ListCredentialsResponse, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *services.ListCredentialsRequest) *services.ListCredentialsResponse); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*services.ListCredentialsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *services.ListCredentialsRequest) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockListCredentials_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockListCredentials_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - data *services.ListCredentialsRequest
func (_e *MockListCredentials_Expecter) Exec(ctx interface{}, data interface{}) *MockListCredentials_Exec_Call {
	return &MockListCredentials_Exec_Call{Call: _e.mock.On("Exec", ctx, data)}
}

func (_c *MockListCredentials_Exec_Call) Run(run func(ctx context.Context, data *services.ListCredentialsRequest)) *MockListCredentials_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*services.ListCredentialsRequest))
	})
	return _c
}

func (_c *MockListCredentials_Exec_Call) Return(_a0 *services.ListCredentialsResponse, _a1 error) *MockListCredentials_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockListCredentials_Exec_Call) RunAndReturn(run func(context.Context, *services.ListCredentialsRequest) (*services.ListCredentialsResponse, error)) *MockListCredentials_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockListCredentials creates a new instance of MockListCredentials. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockListCredentials(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockListCredentials {
	mock := &MockListCredentials{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
