package account

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type AccountStore interface {
	Insert(ctx context.Context, a *Account) error
	Update(ctx context.Context, a *Account) error
	Get(ctx context.Context, id string) (*Account, error)
	Delete(ctx context.Context, id string) error
}

type Store struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*Store, error) {
	return &Store{db: database}, nil
}

func (store *Store) Insert(ctx context.Context, a *Account) error {
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO account (account_name, account_description) VALUES ($1, $2) RETURNING id, account_name, account_description;`,
		a.Name, a.Description).Scan(&a.ID, &a.Name, &a.Description))
	return err
}

func (store *Store) Update(ctx context.Context, a *Account) error {
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE account SET account_name = $2, account_description = $3 WHERE id = $1 RETURNING id, account_name, account_description;`,
		a.ID, a.Name, a.Description).Scan(&a.ID, &a.Name, &a.Description))
	return err
}

func (store *Store) Get(ctx context.Context, id string) (*Account, error) {
	a := &Account{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_name, account_description FROM account WHERE id = $1;`,
		id).Scan(&a.ID, &a.Name, &a.Description))
	return a, err
}

func (store *Store) Delete(ctx context.Context, id string) error {
	_, err := store.db.Exec(ctx, `DELETE FROM account WHERE id = $1;`, id)
	err = db.PgError(err)
	return err
}
