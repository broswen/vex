package project

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type ProjectStore interface {
	List(ctx context.Context, accountId string, limit, offset int64) ([]*Project, error)
	Insert(ctx context.Context, p *Project) (*Project, error)
	Update(ctx context.Context, p *Project) (*Project, error)
	Get(ctx context.Context, projectId string) (*Project, error)
	Delete(ctx context.Context, projectId string) error
}

type Store struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*Store, error) {
	return &Store{db: database}, nil
}

func (store *Store) List(ctx context.Context, accountId string, limit, offset int64) ([]*Project, error) {
	rows, err := store.db.Query(ctx, `SELECT id, account_id, project_name, project_description, created_on, modified_on FROM project WHERE account_id = $1 OFFSET $2 LIMIT $3;`, accountId, offset, limit)
	err = db.PgError(err)
	if err != nil {
		switch err {
		case db.ErrNotFound:
			return nil, ErrProjectNotFound{err.Error()}
		default:
			return nil, ErrUnknown{err}
		}
	}
	defer rows.Close()
	ps := make([]*Project, 0)
	for rows.Next() {
		p := &Project{}
		err = rows.Scan(&p.ID, &p.AccountID, &p.Name, &p.Description, &p.CreatedOn, &p.ModifiedOn)
		if err != nil {
			return nil, ErrUnknown{err}
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (store *Store) Insert(ctx context.Context, p *Project) (*Project, error) {
	newProject := &Project{}
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO project (account_id, project_name, project_description) VALUES ($1, $2, $3) RETURNING id, account_id, project_name, project_description, created_on, modified_on;`,
		p.AccountID, p.Name, p.Description).Scan(&newProject.ID, &newProject.AccountID, &newProject.Name, &newProject.Description, &newProject.CreatedOn, &newProject.ModifiedOn))

	if err != nil {
		switch err {
		case db.ErrNotFound:
			return newProject, ErrProjectNotFound{err.Error()}
		case db.ErrInvalidData:
			return newProject, ErrInvalidData{err.Error()}
		default:
			return newProject, ErrUnknown{err}
		}
	}

	return newProject, nil
}

func (store *Store) Update(ctx context.Context, p *Project) (*Project, error) {
	newProject := &Project{}
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE project SET project_name = $2, project_description = $3 WHERE id = $1 RETURNING id, account_id, project_name, project_description, created_on, modified_on;`,
		p.ID, p.Name, p.Description).Scan(&newProject.ID, &newProject.AccountID, &newProject.Name, &newProject.Description, &newProject.CreatedOn, &newProject.ModifiedOn))

	if err != nil {
		switch err {
		case db.ErrNotFound:
			return newProject, ErrProjectNotFound{err.Error()}
		case db.ErrInvalidData:
			return newProject, ErrInvalidData{err.Error()}
		default:
			return newProject, ErrUnknown{err}
		}
	}
	return newProject, nil
}

func (store *Store) Get(ctx context.Context, projectId string) (*Project, error) {
	p := &Project{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_id, project_name, project_description, created_on, modified_on FROM project WHERE id = $1;`,
		projectId).Scan(&p.ID, &p.AccountID, &p.Name, &p.Description, &p.CreatedOn, &p.ModifiedOn))

	if err != nil {
		switch err {
		case db.ErrNotFound:
			return p, ErrProjectNotFound{err.Error()}
		case db.ErrInvalidData:
			return p, ErrInvalidData{err.Error()}
		default:
			return p, ErrUnknown{err}
		}
	}
	return p, nil
}

func (store *Store) Delete(ctx context.Context, projectId string) error {
	res, err := store.db.Exec(ctx, `DELETE FROM project WHERE id = $1;`, projectId)
	err = db.PgError(err)
	if res.RowsAffected() == 0 && err == nil {
		return ErrProjectNotFound{db.ErrNotFound.Error()}
	}

	if err != nil {
		switch err {
		case db.ErrNotFound:
			return ErrProjectNotFound{err.Error()}
		case db.ErrInvalidData:
			return ErrInvalidData{err.Error()}
		default:
			return ErrUnknown{err}
		}
	}
	return nil
}
