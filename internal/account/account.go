package account

import "time"

type Account struct {
	ID          string    `json:"id"`
	Name        string    `json:"name" db:"account_name"`
	Description string    `json:"description" db:"account_description"`
	CreatedOn   time.Time `json:"created_on" db:"created_on"`
	ModifiedOn  time.Time `json:"modified_on" db:"modified_on"`
}
