package repositoryImplementation

import (
	"errors"

	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/repository"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
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
		return nil,  result.Error
	}
	return &user, nil
}

func (r *userRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}


func (r *userRepository) Update(user *model.User) error {
	return r.db.Transaction(func (tx *gorm.DB) error  {
		if err := tx.Save(user).Error; err != nil {
			return err
		}
		return nil
	}) 
	
}

func (r *userRepository) Delete(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.User{}, "id = ?", id).Error; err != nil {
			return err
		}
		return nil
	})
}

func (r *userRepository) IsEmailUnique(email string) (bool, error) {
	var emailCount int64
	err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&emailCount).Error
	if err != nil {
		return false, err
	}
	return emailCount > 0, nil
}