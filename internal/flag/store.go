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
	rows, err := store.db.Query(ctx, `SELECT id, flag_key, flag_type, flag_value, project_id, account_id, created_on, modified_on FROM flag WHERE project_id = $1 OFFSET $2 LIMIT $3;`, projectId, offset, limit)
	err = db.PgError(err)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return nil, ErrFlagNotFound{err}
		case db.ErrInvalidData:
			return nil, ErrInvalidData{err}
		default:
			return nil, ErrUnknown{err}
		}
	}
	defer rows.Close()
	fs := make([]*Flag, 0)
	for rows.Next() {
		f := &Flag{}
		err = rows.Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID, &f.CreatedOn, &f.ModifiedOn)
		if err != nil {
			return nil, ErrUnknown{err}
		}
		fs = append(fs, f)
	}
	return fs, nil
}

func (store *Store) Insert(ctx context.Context, f *Flag) error {
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO flag (flag_key, flag_type, flag_value, project_id, account_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, flag_key, flag_type, flag_value, project_id, account_id, created_on, modified_on;`,
		f.Key, f.Type, f.Value, f.ProjectID, f.AccountID).Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID, &f.CreatedOn, &f.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return ErrFlagNotFound{err}
		case db.ErrKeyNotUnique:
			return ErrKeyNotUnique{err}
		case db.ErrInvalidData:
			return ErrInvalidData{err}
		default:
			return ErrUnknown{err}
		}
	}
	return nil
}

func (store *Store) Update(ctx context.Context, f *Flag) error {
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE flag SET flag_key = $2, flag_type = $3, flag_value = $4, project_id = $5, account_id = $6 WHERE id = $1 RETURNING id, flag_key, flag_type, flag_value, project_id, account_id, created_on, modified_on;`,
		f.ID, f.Key, f.Type, f.Value, f.ProjectID, f.AccountID).Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID, &f.CreatedOn, &f.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return ErrFlagNotFound{err}
		case db.ErrKeyNotUnique:
			return ErrKeyNotUnique{err}
		case db.ErrInvalidData:
			return ErrInvalidData{err}
		default:
			return ErrUnknown{err}
		}
	}
	return nil
}

func (store *Store) Get(ctx context.Context, id string) (*Flag, error) {
	f := &Flag{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, flag_key, flag_type, flag_value, project_id, account_id, created_on, modified_on FROM flag WHERE id = $1;`,
		id).Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID, &f.CreatedOn, &f.ModifiedOn))

	if err != nil {
		switch err {
		case db.ErrNotFound:
			return f, ErrFlagNotFound{err}
		case db.ErrKeyNotUnique:
			return f, ErrKeyNotUnique{err}
		case db.ErrInvalidData:
			return f, ErrInvalidData{err}
		default:
			return f, ErrUnknown{err}
		}
	}
	return f, nil
}

func (store *Store) Delete(ctx context.Context, id string) error {
	res, err := store.db.Exec(ctx, `DELETE FROM flag WHERE id = $1;`, id)
	err = db.PgError(err)
	if res.RowsAffected() == 0 && err == nil {
		return ErrFlagNotFound{db.ErrNotFound}
	}

	if err != nil {
		switch err {
		case db.ErrNotFound:
			return ErrFlagNotFound{err}
		case db.ErrKeyNotUnique:
			return ErrKeyNotUnique{err}
		case db.ErrInvalidData:
			return ErrInvalidData{err}
		default:
			return ErrUnknown{err}
		}
	}
	return nil
}
