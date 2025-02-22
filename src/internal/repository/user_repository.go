package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/hafiztri123/src/internal/model"
	"gorm.io/gorm"
)

type UserRepository interface {
    GetByID(id string) (*model.User, error)
    GetByEmail(email string) (*model.User, error)
    IsEmailUnique(email string) (bool, error)
    Create(user *model.User) error
    Update(id string, updatedUser *model.User) error
    Delete(id string) error
	ChangePassword(id string, password string) error
}

type userRepository struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
    return &userRepository{
        db: db,
    }
}

func (r *userRepository) GetByID(id string) (*model.User, error) {
    var user model.User
    result := r.db.Where("id = ?", id).First(&user)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, result.Error
    }
    return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
    var user model.User
    result := r.db.Where("email = ?", email).First(&user)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, result.Error
    }
    return &user, nil
}

func (r *userRepository) IsEmailUnique(email string) (bool, error) {
    var emailCount int64
    err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&emailCount).Error
    if err != nil {
        return false, err
    }
    return emailCount == 0, nil
}

func (r *userRepository) Create(user *model.User) error {
    return r.db.Create(user).Error
}

func (r *userRepository) Update(id string, updatedUser *model.User) error {
	var existingUser model.User
	
	userModel := r.db.Model(&model.User{})
	err := userModel.Where("id = ?", id).First(&existingUser).Error
	if err != nil {
		return fmt.Errorf("[FAIL] failed to get user")
	}

	if updatedUser.FullName != "" {
		existingUser.FullName = updatedUser.FullName
	}

	if updatedUser.PhoneNumber != "" {
		existingUser.PhoneNumber = updatedUser.PhoneNumber
	}

	if updatedUser.Organization != "" {
		existingUser.Organization = updatedUser.Organization
	}

	if updatedUser.Bio != "" {
		existingUser.Bio = updatedUser.Bio
	}

	if updatedUser.FullName != "" || updatedUser.PhoneNumber != "" || updatedUser.Organization != "" ||  updatedUser.Bio != "" {
		existingUser.UpdatedAt = time.Now()
	}

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(existingUser).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("[FAIL] update user profile failed: %v", err)
	}

	return nil
}

func (r *userRepository) Delete(id string) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Delete(&model.User{}, "id = ?", id).Error; err != nil {
            return err
        }
        return nil
    })
}

func (r *userRepository) ChangePassword(id string, password string) error {
	userModel := r.db.Model(&model.User{})
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := userModel.Where("id = ?", id).Update("password", password).Error
		if err != nil {
			return err
		}
		return nil
	}) 
	if err != nil {
		return fmt.Errorf("[FAIL] failed to change password: %v", err)
	}
	return nil
}