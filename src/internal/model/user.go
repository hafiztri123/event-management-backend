package model

import "time"

type User struct {
	ID 				string 		`gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email 			string 		`gorm:"type:varchar(255);unique;not null" json:"email"`
	Password 		string 		`gorm:"type:varchar(255);not null" json:"-"`
	FullName 		string 		`gorm:"type:varchar(255);not null" json:"full_name"`
	Role 			string 		`gorm:"type:varchar(20);not null;default:'user'" json:"role"`
	ProfileImage 	string 		`gorm:"type:text" json:"profile_image"`
	PhoneNumber 	string 		`gorm:"type:varchar(20)" json:"phone_number"`
	Organization 	string 		`gorm:"type:varchar(255)" json:"organization"`
	Bio 			string 		`gorm:"type:text" json:"bio"` 
	CreatedAt 		time.Time 	`gorm:"not null" json:"created_at"`
	UpdatedAt 		time.Time 	`gorm:"not null" json:"updated_at"`
	LastLoginAt		*time.Time 	`json:"last_login_at"`
}




type UpdateProfileInput struct {
    FullName     string `json:"full_name"`
    PhoneNumber  string `json:"phone_number"`
    Organization string `json:"organization"`
    Bio         string `json:"bio"`
}

type ChangePasswordInput struct {
    CurrentPassword string `json:"current_password" validate:"required"`
    NewPassword     string `json:"new_password" validate:"required,min=6"`
}

type FileInput struct {
    FileName string
    FileType string
    FileData []byte
}