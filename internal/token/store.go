package token

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"github.com/broswen/vex/internal/db"
)

type Store interface {
	List(ctx context.Context, accountId string, limit, offset int64) ([]*Token, error)
	Generate(ctx context.Context, accountId string, readOnly bool) (*Token, error)
	Reroll(ctx context.Context, tokenId string) (*Token, error)
	Get(ctx context.Context, id string) (*Token, error)
	GetByToken(ctx context.Context, token string) (*Token, error)
	Delete(ctx context.Context, id string) error
}

type PostgresStore struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*PostgresStore, error) {
	return &PostgresStore{db: database}, nil
}

func GenerateTokenAndHash(length int) (string, []byte, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", nil, err
	}
	token := hex.EncodeToString(b)
	hasher := sha256.New()
	hasher.Write([]byte(token))
	hash := hasher.Sum(nil)
	return token, hash, nil
}

func (store *PostgresStore) Generate(ctx context.Context, accountId string, readOnly bool) (*Token, error) {
	t := &Token{}
	generatedToken, tokenHash, err := GenerateTokenAndHash(16)
	if err != nil {
		return nil, err
	}
	err = db.PgError(store.db.QueryRow(ctx, `INSERT INTO token (account_id, read_only, token_hash) VALUES ($1, $2, $3) RETURNING id, account_id, read_only, created_on, modified_on;`,
		accountId, readOnly, tokenHash).Scan(&t.ID, &t.AccountID, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn))
	t.Token = generatedToken
	return t, err
}

func (store *PostgresStore) Reroll(ctx context.Context, tokenId string) (*Token, error) {
	generatedToken, tokenHash, err := GenerateTokenAndHash(16)
	if err != nil {
		return nil, err
	}
	updatedToken := &Token{}
	err = db.PgError(store.db.QueryRow(ctx, `UPDATE token SET token_hash = $1 WHERE id = $2 RETURNING id, account_id, read_only, created_on, modified_on;`,
		tokenHash, tokenId).Scan(&updatedToken.ID, &updatedToken.AccountID, &updatedToken.ReadOnly, &updatedToken.CreatedOn, &updatedToken.ModifiedOn))
	updatedToken.Token = generatedToken
	return updatedToken, err
}

func (store *PostgresStore) Get(ctx context.Context, id string) (*Token, error) {
	t := &Token{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_id, token_hash, read_only, created_on, modified_on FROM token WHERE id = $1;`,
		id).Scan(&t.ID, &t.AccountID, &t.TokenHash, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn))
	return t, err
}

func (store *PostgresStore) GetByToken(ctx context.Context, token string) (*Token, error) {
	t := &Token{}
	hasher := sha256.New()
	hasher.Write([]byte(token))
	hash := hasher.Sum(nil)
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_id, token_hash, read_only, created_on, modified_on FROM token WHERE token_hash = $1;`,
		hash).Scan(&t.ID, &t.AccountID, &t.TokenHash, &t.ReadOnly, &t.CreatedOn, &t.ModifiedOn))
	return t, err
}

func (store *PostgresStore) List(ctx context.Context, accountId string, limit, offset int64) ([]*Token, error) {
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

func (store *PostgresStore) Delete(ctx context.Context, id string) error {
	_, err := store.db.Exec(ctx, `DELETE FROM token WHERE id = $1;`, id)
	err = db.PgError(err)
	return err
}
