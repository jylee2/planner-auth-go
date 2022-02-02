package models

type User struct {
	Uuid     uint
	Name     string
	Email    string
	Password []byte
}