package account

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type Store interface {
	Insert(ctx context.Context, a *Account) (*Account, error)
	Update(ctx context.Context, a *Account) (*Account, error)
	Get(ctx context.Context, id string) (*Account, error)
	List(ctx context.Context, limit, offset int64) ([]*Account, error)
	Delete(ctx context.Context, id string) error
}

type PostgresStore struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*PostgresStore, error) {
	return &PostgresStore{db: database}, nil
}

func (store *PostgresStore) Insert(ctx context.Context, a *Account) (*Account, error) {
	newAccount := &Account{}
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO account (account_name, account_description) VALUES ($1, $2) RETURNING id, account_name, account_description, created_on, modified_on;`,
		a.Name, a.Description).Scan(&newAccount.ID, &newAccount.Name, &newAccount.Description, &newAccount.CreatedOn, &newAccount.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return newAccount, ErrAccountNotFound{err}
		case db.ErrInvalidData:
			return newAccount, ErrInvalidData{err}
		default:
			return newAccount, ErrUnknown{err}
		}
	}
	return newAccount, nil
}

func (store *PostgresStore) Update(ctx context.Context, a *Account) (*Account, error) {
	newAccount := &Account{}
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE account SET account_name = $2, account_description = $3 WHERE id = $1 RETURNING id, account_name, account_description, created_on, modified_on;`,
		a.ID, a.Name, a.Description).Scan(&newAccount.ID, &newAccount.Name, &newAccount.Description, &newAccount.CreatedOn, &newAccount.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return newAccount, ErrAccountNotFound{err}
		case db.ErrInvalidData:
			return newAccount, ErrInvalidData{err}
		default:
			return newAccount, ErrUnknown{err}
		}
	}
	return newAccount, nil
}

func (store *PostgresStore) Get(ctx context.Context, id string) (*Account, error) {
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

func (store *PostgresStore) List(ctx context.Context, limit, offset int64) ([]*Account, error) {
	rows, err := store.db.Query(ctx, `SELECT id, account_name, account_description, created_on, modified_on FROM account LIMIT $1 OFFSET $2;`, limit, offset)
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

func (store *PostgresStore) Delete(ctx context.Context, id string) error {
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
