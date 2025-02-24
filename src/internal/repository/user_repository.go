package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"github.com/jackc/pgx/v5/pgconn"
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
	ChangePhotoProfile(id string, imageURL string) error
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
    err := r.db.Where("id = ?", id).First(&user).Error

	if err != nil {
		return nil, DBError(err)
	}

    return &user, nil
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
    var user model.User
    err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, DBError(err)
	}

    return &user, nil
}

func (r *userRepository) IsEmailUnique(email string) (bool, error) {
    var emailCount int64
    err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&emailCount).Error
    if err != nil {
        return false, DBError(err)
    }
    return emailCount == 0, nil
}

func (r *userRepository) Create(user *model.User) error {
    err := r.db.Create(user).Error
	return DBError(err)
}

func (r *userRepository) Update(id string, updatedUser *model.User) error {
	var existingUser model.User
	
	userModel := r.db.Model(&model.User{})
	err := userModel.Where("id = ?", id).First(&existingUser).Error
	
	if err != nil {
		return DBError(err)
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
		return DBError(err)
	}

	return nil
}

func (r *userRepository) Delete(id string) error {
    return DBError(r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Delete(&model.User{}, "id = ?", id).Error; err != nil {
            return err
        }
        return nil
    }))
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
		return DBError(err)
	}
	return nil

}

func (r *userRepository) ChangePhotoProfile(id string, PhotoProfileURL string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Model(&model.User{}).Where("id = ?", id).Update("profile_image", PhotoProfileURL).Error
		return err
	})

	if err != nil {
		return DBError(err)
	}

	return nil
}





func DBError(err error) error {
    if err == nil {
        return nil // No error to handle
    }

    switch {
    case errors.Is(err, gorm.ErrRecordNotFound):
        return errs.NewNotFoundError("Record not found")

    case errors.Is(err, gorm.ErrInvalidDB):
        return errs.NewValidationError("Invalid database operation")

    case errors.Is(err, gorm.ErrInvalidTransaction):
        return errs.NewValidationError("Invalid transaction")

    case errors.Is(err, gorm.ErrMissingWhereClause):
        return errs.NewValidationError("Missing WHERE clause in query")

    default:
        // Check for PostgreSQL-specific errors
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) {
            switch pgErr.Code {
            case "23505": // Unique violation
                return errs.NewDuplicateEntryError("Duplicate entry: a record with this value already exists")

            case "23503": // Foreign key violation
                return errs.NewValidationError("Foreign key violation: referenced record does not exist")

            case "23514": // Check violation
                return errs.NewValidationError("Check constraint violation: invalid data")

            default:
                return errs.NewDatabaseError(fmt.Sprintf("PostgreSQL error (%s): %s", pgErr.Code, pgErr.Message))
            }
        }

        // Handle general database errors
        return errs.NewDatabaseError(fmt.Sprintf("An unexpected database error occurred: %v", err))
    }
}