package project

import "time"

type Project struct {
	ID          string    `json:"id"`
	AccountID   string    `json:"account_id" db:"account_id"`
	Name        string    `json:"name" db:"project_name"`
	Description string    `json:"description" db:"project_description"`
	CreatedOn   time.Time `json:"created_on" db:"created_on"`
	ModifiedOn  time.Time `json:"modified_on" db:"modified_on"`
}
