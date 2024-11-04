// Code generated by mockery v2.46.3. DO NOT EDIT.

package daomocks

import (
	context "context"

	entities "github.com/a-novel/uservice-credentials/pkg/entities"

	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
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

// Exec provides a mock function with given fields: ctx, ids
func (_m *MockListCredentials) Exec(ctx context.Context, ids []uuid.UUID) ([]*entities.Credential, error) {
	ret := _m.Called(ctx, ids)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 []*entities.Credential
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) ([]*entities.Credential, error)); ok {
		return rf(ctx, ids)
	}
	if rf, ok := ret.Get(0).(func(context.Context, []uuid.UUID) []*entities.Credential); ok {
		r0 = rf(ctx, ids)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Credential)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, []uuid.UUID) error); ok {
		r1 = rf(ctx, ids)
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
//   - ids []uuid.UUID
func (_e *MockListCredentials_Expecter) Exec(ctx interface{}, ids interface{}) *MockListCredentials_Exec_Call {
	return &MockListCredentials_Exec_Call{Call: _e.mock.On("Exec", ctx, ids)}
}

func (_c *MockListCredentials_Exec_Call) Run(run func(ctx context.Context, ids []uuid.UUID)) *MockListCredentials_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]uuid.UUID))
	})
	return _c
}

func (_c *MockListCredentials_Exec_Call) Return(_a0 []*entities.Credential, _a1 error) *MockListCredentials_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockListCredentials_Exec_Call) RunAndReturn(run func(context.Context, []uuid.UUID) ([]*entities.Credential, error)) *MockListCredentials_Exec_Call {
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
