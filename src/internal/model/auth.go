package model

type RegisterInput struct {
    Email    string `json:"email" validate:"required,email" example:"user@example.com"`
    Password string `json:"password" validate:"required,min=6" example:"password123"`
    FullName string `json:"full_name" validate:"required" example:"John Doe"`
}

// LoginInput represents the login request payload
type LoginInput struct {
    Email    string `json:"email" validate:"required,email" example:"user@example.com"`
    Password string `json:"password" validate:"required,min=6" example:"password123"`
}

// LoginResponse represents the login response payload
type LoginResponse struct {
    Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIs..."`
}
