package account

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMockAccountStore(t *testing.T) {
	store := NewMockStore()
	a1 := &Account{
		Name:        "test",
		Description: "test account",
		CreatedOn:   time.Time{},
		ModifiedOn:  time.Time{},
	}

	//test create account
	store.On("Insert").Return(nil)

	err := store.Insert(context.Background(), a1)
	assert.Nil(t, err)
	assert.Equalf(t, "0", a1.ID, "expected 0 for first account id")

	//test get account
	store.On("Get").Return(nil)

	a2, err := store.Get(context.Background(), a1.ID)
	assert.Nil(t, err)
	assert.Equalf(t, a1.ID, a2.ID, "account IDs should match")
	assert.Equalf(t, a1.Name, a2.Name, "account namess should match")

	//test update account
	a2.Name = "updated test"

	store.On("Update").Return(nil)

	err = store.Update(context.Background(), a2)
	assert.Nil(t, err)
	assert.Equalf(t, "updated test", a2.Name, "account name should be updated")

	store.On("Delete").Return(nil)
	err = store.Delete(context.Background(), a2.ID)
	assert.Nil(t, err)

	//test get account after delete
	store.On("Get").Return(nil)

	_, err = store.Get(context.Background(), a2.ID)
	assert.NotNil(t, err)
}
