package project

type Project struct {
	ID          string `json:"id"`
	AccountID   string `json:"account_id" db:"account_id"`
	Name        string `json:"name" db:"project_name"`
	Description string `json:"description" db:"project_description"`
}
