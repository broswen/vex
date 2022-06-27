package project

import (
	"context"
	"github.com/broswen/vex/internal/db"
)

type ProjectStore interface {
	List(ctx context.Context, accountId string, limit, offset int64) ([]*Project, error)
	Insert(ctx context.Context, p *Project) error
	Update(ctx context.Context, p *Project) error
	Get(ctx context.Context, projectId, accountId string) (*Project, error)
	Delete(ctx context.Context, projectId, accountId string) error
}

type Store struct {
	db *db.Database
}

func NewPostgresStore(database *db.Database) (*Store, error) {
	return &Store{db: database}, nil
}

func (store *Store) List(ctx context.Context, accountId string, limit, offset int64) ([]*Project, error) {
	rows, err := store.db.Query(ctx, `SELECT id, account_id, project_name, project_description FROM project WHERE account_id = $1 OFFSET $2 LIMIT $3;`, accountId, offset, limit)
	err = db.PgError(err)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	ps := make([]*Project, 0)
	for rows.Next() {
		p := &Project{}
		err = rows.Scan(&p.ID, &p.AccountID, &p.Name, &p.Description)
		if err != nil {
			return nil, err
		}
		ps = append(ps, p)
	}
	return ps, nil
}

func (store *Store) Insert(ctx context.Context, p *Project) error {
	err := db.PgError(store.db.QueryRow(ctx, `INSERT INTO project (account_id, project_name, project_description) VALUES ($1, $2, $3) RETURNING id, account_id, project_name, project_description;`,
		p.AccountID, p.Name, p.Description).Scan(&p.ID, &p.AccountID, &p.Name, &p.Description))
	return err
}

func (store *Store) Update(ctx context.Context, p *Project) error {
	err := db.PgError(store.db.QueryRow(ctx, `UPDATE project SET project_name = $2, project_description = $3 WHERE id = $1 RETURNING id, account_id, project_name, project_description;`,
		p.ID, p.Name, p.Description).Scan(&p.ID, &p.AccountID, &p.Name, &p.Description))
	return err
}

func (store *Store) Get(ctx context.Context, projectId, accountId string) (*Project, error) {
	p := &Project{}
	err := db.PgError(store.db.QueryRow(ctx, `SELECT id, account_id, project_name, project_description FROM project WHERE id = $1 AND account_id = $2;`,
		projectId, accountId).Scan(&p.ID, &p.AccountID, &p.Name, &p.Description))
	return p, err
}

func (store *Store) Delete(ctx context.Context, projectId, accountId string) error {
	res, err := store.db.Exec(ctx, `DELETE FROM project WHERE id = $1 AND account_id = $2;`, projectId, accountId)
	err = db.PgError(err)
	if res.RowsAffected() == 0 && err == nil {
		return db.ErrNotFound
	}
	return err
}
