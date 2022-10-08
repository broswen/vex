package project

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

func (m *MockStore) Insert(ctx context.Context, a *Project) (*Project, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*Project), args.Error(1)
}

func (m *MockStore) Update(ctx context.Context, a *Project) (*Project, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*Project), args.Error(1)
}

func (m *MockStore) Get(ctx context.Context, id string) (*Project, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Project), args.Error(1)
}

func (m *MockStore) List(ctx context.Context, accountId string, limit, offset int64) ([]*Project, error) {
	args := m.Called(ctx, accountId, limit, offset)
	return args.Get(0).([]*Project), args.Error(1)
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
