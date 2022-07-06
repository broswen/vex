package flag

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type FlagStore interface {
	List(ctx context.Context, projectId string, limit, offset int64) ([]*Flag, error)
	Insert(ctx context.Context, f *Flag) error
	Update(ctx context.Context, f *Flag) error
	Get(ctx context.Context, id string) (*Flag, error)
	Delete(ctx context.Context, id string) error
}

type Store struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*Store, error) {
	return &Store{db: database}, nil
}

func (store *Store) List(ctx context.Context, projectId string, limit, offset int64) ([]*Flag, error) {
	rows, err := store.db.Query(ctx, `SELECT id, flag_key, flag_type, flag_value, project_id, account_id FROM flag WHERE project_id = $1 OFFSET $2 LIMIT $3;`, projectId, offset, limit)
	err = db.PgError(err)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	fs := make([]*Flag, 0)
	for rows.Next() {
		f := &Flag{}
		err = rows.Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID)
		if err != nil {
			return nil, err
		}
		fs = append(fs, f)
	}
	return fs, nil
}

func (store *Store) Insert(ctx context.Context, f *Flag) error {
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO flag (flag_key, flag_type, flag_value, project_id, account_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, flag_key, flag_type, flag_value, project_id, account_id;`,
		f.Key, f.Type, f.Value, f.ProjectID, f.AccountID).Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID))
	return err
}

func (store *Store) Update(ctx context.Context, f *Flag) error {
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE flag SET flag_key = $2, flag_type = $3, flag_value = $4, project_id = $5, account_id = $6 WHERE id = $1 RETURNING id, flag_key, flag_type, flag_value, project_id, account_id;`,
		f.ID, f.Key, f.Type, f.Value, f.ProjectID, f.AccountID).Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID))
	return err
}

func (store *Store) Get(ctx context.Context, id string) (*Flag, error) {
	f := &Flag{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, flag_key, flag_type, flag_value, project_id, account_id FROM flag WHERE id = $1;`,
		id).Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID))
	return f, err
}

func (store *Store) Delete(ctx context.Context, id string) error {
	res, err := store.db.Exec(ctx, `DELETE FROM flag WHERE id = $1;`, id)
	err = db.PgError(err)
	if res.RowsAffected() == 0 && err == nil {
		return db.ErrNotFound
	}
	return err
}
