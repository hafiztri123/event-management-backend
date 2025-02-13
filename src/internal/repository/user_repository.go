package repository

import "github.com/hafiztri123/src/internal/model"
type UserRepository interface {
	GetByID(id string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Create (user *model.User) error 
	Update (user *model.User) error
	Delete (id string) error 
}