package token

import "time"

type Token struct {
	ID         string    `json:"id"`
	AccountID  string    `json:"accountId" db:"account_id"`
	Token      string    `json:"token"`
	ReadOnly   bool      `json:"readOnly" db:"read_only"`
	CreatedOn  time.Time `json:"created_on" db:"created_on"`
	ModifiedOn time.Time `json:"modified_on" db:"modified_on"`
}
