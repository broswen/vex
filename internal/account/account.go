package account

type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name" db:"account_name"`
	Description string `json:"description" db:"account_description"`
}
