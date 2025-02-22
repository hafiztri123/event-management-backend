package service

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "github.com/hafiztri123/src/internal/model"
    "github.com/hafiztri123/src/internal/pkg/config"
    "github.com/hafiztri123/src/internal/repository"
    "golang.org/x/crypto/bcrypt"
)

// AuthService defines the interface for authentication-related service operations.
type AuthService interface {
    Register(input *model.RegisterInput) error
    Login(input *model.LoginInput) (*model.LoginResponse, error)
}

// AuthService implements the AuthService interface.
type authService struct {
    userRepo repository.UserRepository
    config   *config.AuthConfig
}

// NewAuthService creates a new instance of AuthService.
func NewAuthService(userRepo repository.UserRepository, config *config.AuthConfig) AuthService {
    return &authService{
        userRepo: userRepo,
        config:   config,
    }
}

var _ AuthService = (*authService)(nil)

// Register registers a new user with the provided input.
func (s *authService) Register(input *model.RegisterInput) error {
    emailExist, err := s.userRepo.IsEmailUnique(input.Email)
    if err != nil {
        return err
    }
    if !emailExist {
        return errors.New("[FAIL] user already exists")
    }

    encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return err
    }

    user := &model.User{
        Email:     input.Email,
        Password:  string(encryptedPassword),
        FullName:  input.FullName,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    err = s.userRepo.Create(user)
    if err != nil {
        return err
    }

    return nil
}

// Login authenticates a user and generates a JWT token if credentials are valid.
func (s *authService) Login(input *model.LoginInput) (*model.LoginResponse, error) {
    existingUser, err := s.userRepo.GetByEmail(input.Email)
    if err != nil {
        return nil, err
    }
    if existingUser == nil {
        return nil, errors.New("[FAIL] user not found")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(input.Password)); err != nil {
        return nil, errors.New("[FAIL] invalid credentials")
    }

    token, err := s.generateToken(existingUser)
    if err != nil {
        return nil, err
    }

    return &model.LoginResponse{
        Token: token,
    }, nil
}

// generateToken generates a signed JWT token for the given user.
func (s *authService) generateToken(user *model.User) (string, error) {
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "email":   user.Email,
        "nbf":     time.Now().Unix(),
        "iat":     time.Now().Unix(),
        "exp":     time.Now().Add(time.Hour * time.Duration(s.config.TokenExpiryHrs)).Unix(),
        "role":    user.Role,
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString([]byte(s.config.JWTSecret))
    if err != nil {
        return "", err
    }

    return signedToken, nil
}