package models

type User struct {
	Uuid     string `json:"uuid"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password []byte `json:"-"` // don't return password
}
