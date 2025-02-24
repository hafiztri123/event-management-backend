package service

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"time"

	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"github.com/hafiztri123/src/internal/pkg/storage"
	"github.com/hafiztri123/src/internal/repository"
	"golang.org/x/crypto/bcrypt"
)


type UserService interface {
    UpdateProfile(userID string, input *model.UpdateProfileInput) error
    GetProfile(userID string) (*model.User, error)
    ChangePassword(userID string, input *model.ChangePasswordInput) error
    UploadProfileImage(ctx context.Context, userID string, file multipart.File, fileName string) error
}


type userServiceImpl struct {
	userRepo repository.UserRepository
	storage storage.StorageService
}

func NewUserService(userRepo repository.UserRepository, storage storage.StorageService) UserService {
	return userServiceImpl{
		userRepo: userRepo,
		storage: storage,
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

func (s userServiceImpl) UploadProfileImage(
	 ctx context.Context,
	 userID string, 
	 file multipart.File, 
	 filename string,
	 ) error {

	imageURL, err := s.storage.UploadFile(ctx, file, filename)
	if err != nil {
		return err
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}

	if user.ProfileImage != "" {
		if err := s.storage.DeleteFile(ctx, storage.ExtractPublicID(user.ProfileImage)); err != nil {
			log.Printf("Failed to delete old image: %v", err)
		}
	}

	user.ProfileImage = imageURL
	return s.userRepo.ChangePhotoProfile(userID, imageURL)
	
}


func generateHashPassword(password string) (string, error)  {
	hashedPassword, err :=  bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("password hashing failed: %v", err)
	}

	return string(hashedPassword), nil
}