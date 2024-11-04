// Code generated by mockery v2.46.3. DO NOT EDIT.

package servicesmocks

import (
	context "context"

	services "github.com/a-novel/uservice-credentials/pkg/services"
	mock "github.com/stretchr/testify/mock"
)

// MockGetCredentials is an autogenerated mock type for the GetCredentials type
type MockGetCredentials struct {
	mock.Mock
}

type MockGetCredentials_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGetCredentials) EXPECT() *MockGetCredentials_Expecter {
	return &MockGetCredentials_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, data
func (_m *MockGetCredentials) Exec(ctx context.Context, data *services.GetCredentialsRequest) (*services.GetCredentialsResponse, error) {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 *services.GetCredentialsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *services.GetCredentialsRequest) (*services.GetCredentialsResponse, error)); ok {
		return rf(ctx, data)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *services.GetCredentialsRequest) *services.GetCredentialsResponse); ok {
		r0 = rf(ctx, data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*services.GetCredentialsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *services.GetCredentialsRequest) error); ok {
		r1 = rf(ctx, data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGetCredentials_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockGetCredentials_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - data *services.GetCredentialsRequest
func (_e *MockGetCredentials_Expecter) Exec(ctx interface{}, data interface{}) *MockGetCredentials_Exec_Call {
	return &MockGetCredentials_Exec_Call{Call: _e.mock.On("Exec", ctx, data)}
}

func (_c *MockGetCredentials_Exec_Call) Run(run func(ctx context.Context, data *services.GetCredentialsRequest)) *MockGetCredentials_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*services.GetCredentialsRequest))
	})
	return _c
}

func (_c *MockGetCredentials_Exec_Call) Return(_a0 *services.GetCredentialsResponse, _a1 error) *MockGetCredentials_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockGetCredentials_Exec_Call) RunAndReturn(run func(context.Context, *services.GetCredentialsRequest) (*services.GetCredentialsResponse, error)) *MockGetCredentials_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockGetCredentials creates a new instance of MockGetCredentials. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockGetCredentials(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockGetCredentials {
	mock := &MockGetCredentials{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
