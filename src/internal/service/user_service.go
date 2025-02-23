package service

import (
	"fmt"
	"time"

	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"github.com/hafiztri123/src/internal/repository"
	"golang.org/x/crypto/bcrypt"
)


type UserService interface {
    UpdateProfile(userID string, input *model.UpdateProfileInput) error
    GetProfile(userID string) (*model.User, error)
    ChangePassword(userID string, input *model.ChangePasswordInput) error
    // UploadProfileImage(userID string, fileInput *model.FileInput) error
}


type userServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return userServiceImpl{
		userRepo: userRepo,
	}
}

func (s userServiceImpl) UpdateProfile(userID string, input *model.UpdateProfileInput) error {
	if input == nil {
		return errs.NewBadRequestError("Request is missing")
	}

	updatedModel := &model.User{
		ID: userID,
		FullName: input.FullName,
		PhoneNumber: input.PhoneNumber,
		Organization: input.Organization,
		Bio: input.Bio,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.userRepo.Update(userID,updatedModel); err != nil {
		return err
	}

	return nil
}

func (s userServiceImpl) GetProfile(userID string) (*model.User, error) {
	if userID == "" {
		return nil, errs.NewBadRequestError("Request is missing")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s userServiceImpl)  ChangePassword(userID string, input *model.ChangePasswordInput) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.CurrentPassword)); err != nil {
		return errs.NewForbiddenError("Invalid credentials")

	}
	currentHashedPassword, err := generateHashPassword(input.CurrentPassword)
	if err != nil {
		return errs.NewInternalServerError(err.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(currentHashedPassword), []byte(input.NewPassword)); err == nil {
		return errs.NewDuplicateEntryError("Password cannot be the same as the old password")
	}

	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)

	if err != nil {
		return errs.NewInternalServerError(err.Error())
	}

	return s.userRepo.ChangePassword(userID, string(newHashedPassword))
	
}

// func (s userServiceImpl) UploadProfileImage(ctx context.Context, userID string, fileInput *model.FileInput, filename string) error {
// }


func generateHashPassword(password string) (string, error)  {
	hashedPassword, err :=  bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("Password hashing failed: %v", err)
	}

	return string(hashedPassword), nil
}