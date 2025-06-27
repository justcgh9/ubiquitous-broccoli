package user

type User struct {
    Id      string
    Email   string
    Handle  string
    Token   string
}

type LoginDTO struct {
    Email       string
    Password    string
}

type RegisterDTO struct {
    Handle string
    Email string
    Password string
}

func New(id, email, handle string) User {
    return User{
        Id: id,
        Email: email,
        Handle: handle,
    }
}

func NewLoginDTO(email, password string) LoginDTO {
    return LoginDTO{
        Email: email,
        Password: password,
    }
}

func NewRegisterDTO(
    handle string,
    email string,
    password string,
) RegisterDTO {
    return RegisterDTO{
        Handle: handle,
        Email: email,
        Password: password,
    }
}
