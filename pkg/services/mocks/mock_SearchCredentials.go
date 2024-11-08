// Code generated by mockery v2.46.3. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	services "github.com/a-novel/uservice-credentials/pkg/services"
	mock "github.com/stretchr/testify/mock"
)

// MockSearchCredentials is an autogenerated mock type for the SearchCredentials type
type MockSearchCredentials struct {
	mock.Mock
}

type MockSearchCredentials_Expecter struct {
	mock *mock.Mock
}

func (_m *MockSearchCredentials) EXPECT() *MockSearchCredentials_Expecter {
	return &MockSearchCredentials_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, data
func (_m *MockSearchCredentials) Exec(ctx context.Context, data *services.SearchCredentialsRequest) (*services.SearchCredentialsResponse, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 *services.SearchCredentialsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *services.SearchCredentialsRequest) (*services.SearchCredentialsResponse, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *services.SearchCredentialsRequest) *services.SearchCredentialsResponse); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*services.SearchCredentialsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *services.SearchCredentialsRequest) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockSearchCredentials_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockSearchCredentials_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - data *services.SearchCredentialsRequest
func (_e *MockSearchCredentials_Expecter) Exec(ctx interface{}, data interface{}) *MockSearchCredentials_Exec_Call {
	return &MockSearchCredentials_Exec_Call{Call: _e.mock.On("Exec", ctx, data)}
}

func (_c *MockSearchCredentials_Exec_Call) Run(run func(ctx context.Context, data *services.SearchCredentialsRequest)) *MockSearchCredentials_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*services.SearchCredentialsRequest))
	})
	return _c
}

func (_c *MockSearchCredentials_Exec_Call) Return(_a0 *services.SearchCredentialsResponse, _a1 error) *MockSearchCredentials_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockSearchCredentials_Exec_Call) RunAndReturn(run func(context.Context, *services.SearchCredentialsRequest) (*services.SearchCredentialsResponse, error)) *MockSearchCredentials_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockSearchCredentials creates a new instance of MockSearchCredentials. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockSearchCredentials(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockSearchCredentials {
	mock := &MockSearchCredentials{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
