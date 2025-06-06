package user

type User struct {
    Email   string
    Handle  string
}

type UserLoginDTO struct {
    Email       string
    Password    string
}

func NewLoginDTO(email, password string) UserLoginDTO {
    return UserLoginDTO{
        Email: email,
        Password: password,
    }
}
