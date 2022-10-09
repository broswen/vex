package flag

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

func (m *MockStore) Insert(ctx context.Context, a *Flag) (*Flag, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*Flag), args.Error(1)
}

func (m *MockStore) Update(ctx context.Context, a *Flag) (*Flag, error) {
	args := m.Called(ctx, a)
	return args.Get(0).(*Flag), args.Error(1)
}

func (m *MockStore) Get(ctx context.Context, id string) (*Flag, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Flag), args.Error(1)
}

func (m *MockStore) List(ctx context.Context, projectId string, limit, offset int64) ([]*Flag, error) {
	args := m.Called(ctx, projectId, limit, offset)
	return args.Get(0).([]*Flag), args.Error(1)
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockStore) ReplaceFlags(ctx context.Context, projectId string, flags []*Flag) ([]*Flag, error) {
	args := m.Called(ctx, projectId, flags)
	return args.Get(0).([]*Flag), args.Error(1)
}
