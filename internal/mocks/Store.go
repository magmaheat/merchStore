// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/magmaheat/merchStore/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// Store is an autogenerated mock type for the Store type
type Store struct {
	mock.Mock
}

// BuyItem provides a mock function with given fields: ctx, nameItem
func (_m *Store) BuyItem(ctx context.Context, nameItem string) error {
	ret := _m.Called(ctx, nameItem)

	if len(ret) == 0 {
		panic("no return value specified for BuyItem")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string) error); ok {
		r0 = rf(ctx, nameItem)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetInfo provides a mock function with given fields: ctx
func (_m *Store) GetInfo(ctx context.Context) (*models.Info, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetInfo")
	}

	var r0 *models.Info
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*models.Info, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *models.Info); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.Info)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SendCoin provides a mock function with given fields: ctx, toUser, amount
func (_m *Store) SendCoin(ctx context.Context, toUser string, amount int) error {
	ret := _m.Called(ctx, toUser, amount)

	if len(ret) == 0 {
		panic("no return value specified for SendCoin")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, int) error); ok {
		r0 = rf(ctx, toUser, amount)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewStore creates a new instance of Store. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewStore(t interface {
	mock.TestingT
	Cleanup(func())
}) *Store {
	mock := &Store{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
