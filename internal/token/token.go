package token

import "time"

type Token struct {
	ID         string    `json:"id"`
	AccountID  string    `json:"account_id" db:"account_id"`
	Token      string    `json:"token,omitempty"`
	TokenHash  []byte    `json:"-" db:"token_hash"`
	ReadOnly   bool      `json:"read_only" db:"read_only"`
	CreatedOn  time.Time `json:"created_on" db:"created_on"`
	ModifiedOn time.Time `json:"modified_on" db:"modified_on"`
}
