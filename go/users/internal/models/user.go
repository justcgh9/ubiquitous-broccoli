package models

type User struct {
	ID       int64
	Email    string
	Handle	 string
	PassHash []byte
}