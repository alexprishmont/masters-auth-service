package models

type User struct {
	UniqueId     string
	Email        string
	PasswordHash []byte
}
