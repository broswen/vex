package account

type Account struct {
	ID   string `json:"id"`
	Name string `json:"name" db:"account_name"`
}
