package models

type User struct {
	ID       int64
	Email    string
	Handle	 string
	PassHash []byte
}

type UserDTO struct {
	ID       int64
	Email    string
	Handle	 string
}

func NewDTOFromUser(usr User) UserDTO {
	return UserDTO{
		ID: usr.ID,
		Email: usr.Email,
		Handle: usr.Handle,
	}
}