package user

type Response struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Email  string  `json:"email"`
	APIKey *string `json:"api_key,omitempty"`
}
