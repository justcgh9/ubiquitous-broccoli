package user

type User struct {
    Id      string
    Email   string
    Handle  string
}

type UserLoginDTO struct {
    Email       string
    Password    string
}

func NewUser(id, email, handle string) User {
    return User{
        Id: id,
        Email: email,
        Handle: handle,
    }
}

func NewLoginDTO(email, password string) UserLoginDTO {
    return UserLoginDTO{
        Email: email,
        Password: password,
    }
}
