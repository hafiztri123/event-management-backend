package model

import "time"

type User struct {
	ID 			string 		`gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Email 		string 		`gorm:"type:varchar(255);unique;not null" json:"email"`
	Password 	string 		`gorm:"type:varchar(255);not null" json:"-"`
	FullName 	string 		`gorm:"type:varchar(255);not null" json:"full_name"`
	CreatedAt 	time.Time 	`gorm:"not null" json:"created_at"`
	UpdatedAt 	time.Time 	`gorm:"not null" json:"updated_at"`
}

