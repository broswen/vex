package account

import (
	"context"
	"fmt"
	"github.com/broswen/vex/internal/db"
	"github.com/stretchr/testify/mock"
	"time"
)

type MockStore struct {
	nextAccountID int
	accounts      map[string]*Account
	mock.Mock
}

func NewMockStore() *MockStore {
	return &MockStore{
		nextAccountID: 0,
		accounts:      make(map[string]*Account),
	}
}

func (m *MockStore) Insert(ctx context.Context, a *Account) error {
	args := m.Called()
	a.ID = fmt.Sprint(m.nextAccountID)
	a.CreatedOn = time.Now()
	a.ModifiedOn = time.Now()
	m.nextAccountID++
	m.accounts[a.ID] = a
	return args.Error(0)
}

func (m *MockStore) Update(ctx context.Context, a *Account) error {
	args := m.Called()
	a.ModifiedOn = time.Now()
	m.accounts[a.ID] = a
	return args.Error(0)
}

func (m *MockStore) Get(ctx context.Context, id string) (*Account, error) {
	args := m.Called()
	a, ok := m.accounts[id]
	if !ok {
		return nil, db.ErrNotFound
	}
	return a, args.Error(0)
}

func (m *MockStore) Delete(ctx context.Context, id string) error {
	args := m.Called()
	delete(m.accounts, id)
	return args.Error(0)
}
