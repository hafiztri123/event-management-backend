package model

import "time"

type Category struct {
	ID 			string 		`gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name 		string 		`gorm:"type:varchar(100);unique;not null" json:"name"`
	Description string 		`gorm:"type:text" json:"description"`
	CreatedAt 	time.Time 	`gorm:"not null" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"not null" json:"updated_at"`
}

type CreateCategoryInput struct {
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
}

type UpdateCategoryInput struct {
    Name        string `json:"name" validate:"required"`
    Description string `json:"description"`
}