package token

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type TokenStore interface {
	List(ctx context.Context, accountId string, limit, offset int64) ([]*Token, error)
	Generate(ctx context.Context, accountId string, readOnly bool) (*Token, error)
	Reroll(ctx context.Context, t *Token) error
	Get(ctx context.Context, id string) (*Token, error)
	GetByToken(ctx context.Context, token string) (*Token, error)
	Delete(ctx context.Context, id string) error
}

type Store struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*Store, error) {
	return &Store{db: database}, nil
}

func (store *Store) Generate(ctx context.Context, accountId string, readOnly bool) (*Token, error) {
	t := &Token{}
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO token (account_id, read_only) VALUES ($1, $2) RETURNING id, token, account_id, read_only, created_on, modified_on;`,
		accountId, readOnly).Scan(&t.ID, &t.Token, &t.AccountID, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn))
	return t, err
}

func (store *Store) Reroll(ctx context.Context, t *Token) error {
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE token SET token = uuid_generate_v4() WHERE id = $1 AND account_id = $2 RETURNING id, token, account_id, read_only, created_on, modified_on;`,
		t.ID, t.AccountID).Scan(&t.ID, &t.Token, &t.AccountID, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn))
	return err
}

func (store *Store) Get(ctx context.Context, id string) (*Token, error) {
	t := &Token{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_id, token, read_only, created_on, modified_on FROM token WHERE id = $1;`,
		id).Scan(&t.ID, &t.AccountID, &t.Token, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn))
	return t, err
}

func (store *Store) GetByToken(ctx context.Context, token string) (*Token, error) {
	t := &Token{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_id, token, read_only, created_on, modified_on FROM token WHERE token = $1;`,
		token).Scan(&t.ID, &t.AccountID, &token, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn))
	return t, err
}

func (store *Store) List(ctx context.Context, accountId string, limit, offset int64) ([]*Token, error) {
	rows, err := store.db.Query(ctx, `SELECT id, account_id, read_only, created_on, modified_on FROM token WHERE account_id = $1 OFFSET $2 LIMIT $3;`, accountId, offset, limit)
	err = db.PgError(err)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fs := make([]*Token, 0)
	for rows.Next() {
		t := &Token{}
		err = rows.Scan(&t.ID, &t.AccountID, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn)
		if err != nil {
			return nil, err
		}
		fs = append(fs, t)
	}
	return fs, nil
}

func (store *Store) Delete(ctx context.Context, id string) error {
	_, err := store.db.Exec(ctx, `DELETE FROM token WHERE id = $1;`, id)
	err = db.PgError(err)
	return err
}
