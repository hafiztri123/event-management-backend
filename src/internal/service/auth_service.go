package service


type AuthService interface {
	Register(input *RegisterInput) error
	Login(input *LoginInput) (*LoginResponse, error)
}

type RegisterInput struct {
	Email 		string 	`json:"email" validate:"required,email"`
	Password 	string 	`json:"password" validate:"required,min=6"`
	FullName 	string 	`json:"full_name" validate:"required"`
}

type LoginInput struct {
	Email 		string 	`json:"email" validate:"required,email"`
	Password 	string 	`json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

