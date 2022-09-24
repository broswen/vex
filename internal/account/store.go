package account

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type AccountStore interface {
	Insert(ctx context.Context, a *Account) error
	Update(ctx context.Context, a *Account) error
	Get(ctx context.Context, id string) (*Account, error)
	List(ctx context.Context) ([]*Account, error)
	Delete(ctx context.Context, id string) error
}

type Store struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*Store, error) {
	return &Store{db: database}, nil
}

func (store *Store) Insert(ctx context.Context, a *Account) error {
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO account (account_name, account_description) VALUES ($1, $2) RETURNING id, account_name, account_description, created_on, modified_on;`,
		a.Name, a.Description).Scan(&a.ID, &a.Name, &a.Description, &a.CreatedOn, &a.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return ErrAccountNotFound{err}
		case db.ErrInvalidData:
			return ErrInvalidData{err}
		default:
			return ErrUnknown{err}
		}
	}
	return nil
}

func (store *Store) Update(ctx context.Context, a *Account) error {
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE account SET account_name = $2, account_description = $3 WHERE id = $1 RETURNING id, account_name, account_description, created_on, modified_on;`,
		a.ID, a.Name, a.Description).Scan(&a.ID, &a.Name, &a.Description, &a.CreatedOn, &a.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return ErrAccountNotFound{err}
		case db.ErrInvalidData:
			return ErrInvalidData{err}
		default:
			return ErrUnknown{err}
		}
	}
	return nil
}

func (store *Store) Get(ctx context.Context, id string) (*Account, error) {
	a := &Account{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_name, account_description, created_on, modified_on FROM account WHERE id = $1;`,
		id).Scan(&a.ID, &a.Name, &a.Description, &a.CreatedOn, &a.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return a, ErrAccountNotFound{err}
		case db.ErrInvalidData:
			return a, ErrInvalidData{err}
		default:
			return a, ErrUnknown{err}
		}
	}
	return a, nil
}

func (store *Store) List(ctx context.Context) ([]*Account, error) {
	rows, err := store.db.Query(ctx, `SELECT id, account_name, account_description, created_on, modified_on FROM account;`)
	err = db.PgError(err)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return nil, ErrAccountNotFound{err}
		case db.ErrInvalidData:
			return nil, ErrInvalidData{err}
		default:
			return nil, ErrUnknown{err}
		}
	}

	defer rows.Close()
	accounts := make([]*Account, 0)
	for rows.Next() {
		a := &Account{}
		err = rows.Scan(&a.ID, &a.Name, &a.Description, &a.CreatedOn, &a.ModifiedOn)
		if err != nil {
			return nil, ErrUnknown{err}
		}
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func (store *Store) Delete(ctx context.Context, id string) error {
	_, err := store.db.Exec(ctx, `DELETE FROM account WHERE id = $1;`, id)
	err = db.PgError(err)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return ErrAccountNotFound{err}
		case db.ErrInvalidData:
			return ErrInvalidData{err}
		default:
			return ErrUnknown{err}
		}
	}
	return nil
}
