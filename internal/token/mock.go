package token

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

func (m *MockStore) Generate(ctx context.Context, accountId string, readOnly bool) (*Token, error) {
	args := m.Called(ctx, accountId, readOnly)
	return args.Get(0).(*Token), args.Error(1)
}

func (m *MockStore) Reroll(ctx context.Context, tokenId string) (*Token, error) {
	args := m.Called(ctx, tokenId)
	return args.Get(0).(*Token), args.Error(1)
}

func (m *MockStore) Get(ctx context.Context, id string) (*Token, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*Token), args.Error(1)
}

func (m *MockStore) GetByToken(ctx context.Context, token string) (*Token, error) {
	args := m.Called(ctx, token)
	return args.Get(0).(*Token), args.Error(1)
}

func (m *MockStore) List(ctx context.Context, accountId string, limit, offset int64) ([]*Token, error) {
	args := m.Called(ctx, accountId, limit, offset)
	return args.Get(0).([]*Token), args.Error(1)
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
