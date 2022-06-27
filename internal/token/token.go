package token

type Token struct {
	ID        string `json:"id"`
	AccountID string `json:"accountId" db:"account_id"`
	Token     string `json:"token"`
	ReadOnly  bool   `json:"readOnly" db:"read_only"`
}
