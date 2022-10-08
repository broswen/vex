package account

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMockAccountStore(t *testing.T) {
	store := NewMockStore()
	ctx := context.Background()
	now := time.Now()
	a1 := &Account{
		Name:        "test",
		Description: "test account",
		CreatedOn:   now,
		ModifiedOn:  now,
	}

	//test create account
	store.On("Insert", ctx, a1).Return(&Account{
		ID:          "0",
		Name:        "test",
		Description: "test account",
		CreatedOn:   now,
		ModifiedOn:  now,
	}, nil)

	a1, err := store.Insert(ctx, a1)
	assert.Nil(t, err)
	assert.Equalf(t, "0", a1.ID, "expected 0 for first account id")

	//test get account
	store.On("Get", ctx, a1.ID).Return(&Account{
		ID:          "0",
		Name:        "test",
		Description: "test account",
		CreatedOn:   now,
		ModifiedOn:  now,
	}, nil)

	a2, err := store.Get(ctx, a1.ID)
	assert.Nil(t, err)
	assert.Equal(t, a1, a2)

	store.On("Delete", ctx, a2.ID).Return(nil)
	err = store.Delete(context.Background(), a2.ID)
	assert.Nil(t, err)
}
