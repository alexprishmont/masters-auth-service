package models

type User struct {
	UniqueId     string       `bson:"uniqueId"`
	Email        string       `bson:"email"`
	PasswordHash []byte       `bson:"-"`
	Permissions  []Permission `bson:"permissions"`
}

type Permission struct {
	Name string `bson:"name"`
}
