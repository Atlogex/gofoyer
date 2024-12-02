package models

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Email     string
	IsAdmin   bool
	PassHash  []byte
}
