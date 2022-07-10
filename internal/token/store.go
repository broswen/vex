package token

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type TokenStore interface {
	Generate(ctx context.Context, accountId string, readOnly bool) (*Token, error)
	Reroll(ctx context.Context, t *Token) error
	Get(ctx context.Context, id string) (*Token, error)
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
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_id, read_only, created_on, modified_on FROM token WHERE token = $1;`,
		id).Scan(&t.ID, &t.AccountID, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn))
	return t, err
}

func (store *Store) Delete(ctx context.Context, id string) error {
	_, err := store.db.Exec(ctx, `DELETE FROM token WHERE id = $1;`, id)
	err = db.PgError(err)
	return err
}
