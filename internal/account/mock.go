package account

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func NewMockStore() *MockStore {
	return &MockStore{}
}

func (m *MockStore) Insert(ctx context.Context, a *Account) (*Account, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *MockStore) Update(ctx context.Context, a *Account) (*Account, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *MockStore) Get(ctx context.Context, id string) (*Account, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Account), args.Error(1)
}

func (m *MockStore) List(ctx context.Context, limit, offset int64) ([]*Account, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*Account), args.Error(1)
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
