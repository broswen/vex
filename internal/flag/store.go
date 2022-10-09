package flag

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type Store interface {
	List(ctx context.Context, projectId string, limit, offset int64) ([]*Flag, error)
	Insert(ctx context.Context, f *Flag) (*Flag, error)
	Update(ctx context.Context, f *Flag) (*Flag, error)
	Get(ctx context.Context, id string) (*Flag, error)
	Delete(ctx context.Context, id string) error
}

type PostgresStore struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*PostgresStore, error) {
	return &PostgresStore{db: database}, nil
}

func (store *PostgresStore) List(ctx context.Context, projectId string, limit, offset int64) ([]*Flag, error) {
	rows, err := store.db.Query(ctx, `SELECT id, flag_key, flag_type, flag_value, project_id, account_id, created_on, modified_on FROM flag WHERE project_id = $1 OFFSET $2 LIMIT $3;`, projectId, offset, limit)
	err = db.PgError(err)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return nil, ErrFlagNotFound{err.Error()}
		case db.ErrInvalidData:
			return nil, ErrInvalidData{err.Error()}
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

func (store *PostgresStore) Insert(ctx context.Context, f *Flag) (*Flag, error) {
	newFlag := &Flag{}
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO flag (flag_key, flag_type, flag_value, project_id, account_id) VALUES ($1, $2, $3, $4, $5) RETURNING id, flag_key, flag_type, flag_value, project_id, account_id, created_on, modified_on;`,
		f.Key, f.Type, f.Value, f.ProjectID, f.AccountID).Scan(&newFlag.ID, &newFlag.Key, &newFlag.Type, &newFlag.Value, &newFlag.ProjectID, &newFlag.AccountID, &newFlag.CreatedOn, &newFlag.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return newFlag, ErrFlagNotFound{err.Error()}
		case db.ErrKeyNotUnique:
			return newFlag, ErrKeyNotUnique{err.Error()}
		case db.ErrInvalidData:
			return newFlag, ErrInvalidData{err.Error()}
		default:
			return newFlag, ErrUnknown{err}
		}
	}
	return newFlag, nil
}

func (store *PostgresStore) Update(ctx context.Context, f *Flag) (*Flag, error) {
	updatedFlag := &Flag{}
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE flag SET flag_key = $2, flag_type = $3, flag_value = $4, project_id = $5, account_id = $6 WHERE id = $1 RETURNING id, flag_key, flag_type, flag_value, project_id, account_id, created_on, modified_on;`,
		f.ID, f.Key, f.Type, f.Value, f.ProjectID, f.AccountID).Scan(&updatedFlag.ID, &updatedFlag.Key, &updatedFlag.Type, &updatedFlag.Value, &updatedFlag.ProjectID, &updatedFlag.AccountID, &updatedFlag.CreatedOn, &updatedFlag.ModifiedOn))
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return updatedFlag, ErrFlagNotFound{err.Error()}
		case db.ErrKeyNotUnique:
			return updatedFlag, ErrKeyNotUnique{err.Error()}
		case db.ErrInvalidData:
			return updatedFlag, ErrInvalidData{err.Error()}
		default:
			return updatedFlag, ErrUnknown{err}
		}
	}
	return updatedFlag, nil
}

func (store *PostgresStore) Get(ctx context.Context, id string) (*Flag, error) {
	f := &Flag{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, flag_key, flag_type, flag_value, project_id, account_id, created_on, modified_on FROM flag WHERE id = $1;`,
		id).Scan(&f.ID, &f.Key, &f.Type, &f.Value, &f.ProjectID, &f.AccountID, &f.CreatedOn, &f.ModifiedOn))

	if err != nil {
		switch err {
		case db.ErrNotFound:
			return f, ErrFlagNotFound{err.Error()}
		case db.ErrKeyNotUnique:
			return f, ErrKeyNotUnique{err.Error()}
		case db.ErrInvalidData:
			return f, ErrInvalidData{err.Error()}
		default:
			return f, ErrUnknown{err}
		}
	}
	return f, nil
}

func (store *PostgresStore) Delete(ctx context.Context, id string) error {
	res, err := store.db.Exec(ctx, `DELETE FROM flag WHERE id = $1;`, id)
	err = db.PgError(err)
	if res.RowsAffected() == 0 && err == nil {
		return ErrFlagNotFound{db.ErrNotFound.Error()}
	}

	if err != nil {
		switch err {
		case db.ErrNotFound:
			return ErrFlagNotFound{err.Error()}
		case db.ErrKeyNotUnique:
			return ErrKeyNotUnique{err.Error()}
		case db.ErrInvalidData:
			return ErrInvalidData{err.Error()}
		default:
			return ErrUnknown{err}
		}
	}
	return nil
}
