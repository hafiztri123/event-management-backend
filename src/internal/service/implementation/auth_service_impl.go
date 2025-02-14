package serviceImplementation

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/pkg/config"
	"github.com/hafiztri123/src/internal/repository"
	"github.com/hafiztri123/src/internal/service"
	"golang.org/x/crypto/bcrypt"
)


type AuthService struct {
	userRepo repository.UserRepository
	config *config.AuthConfig
}

func NewAuthService(userRepo repository.UserRepository, config *config.AuthConfig) service.AuthService {
	return &AuthService{
		userRepo: userRepo,
		config: config,
	}
}

var _ service.AuthService = (*AuthService)(nil)

func (s *AuthService) Register(input *service.RegisterInput) error {
	emailExist, err := s.userRepo.IsEmailUnique(input.Email)
	if err != nil {
		return err
	}

	if emailExist {
		return errors.New("[FAIL] user already exists")
	}

	encryptedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	user := &model.User{
		Email: input.Email,
		Password: string(encryptedPassword),
		FullName: input.FullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Login(input *service.LoginInput) (*service.LoginResponse, error) {
	existingUser, err := s.userRepo.GetByEmail(input.Email)

	if err != nil {
		return nil, err
	}

	if  existingUser == nil {
		return nil, errors.New("[FAIL] user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("[FAIL] invalid credentials")
	}

	token, err := s.generateToken(existingUser)
	if err != nil {
		return nil, err
	}

	return &service.LoginResponse{
		Token: token,
	}, nil
}




func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * time.Duration(s.config.TokenExpiryHrs)).Unix(),
	}

	

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}