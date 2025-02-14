package model

import "time"

type Event struct {
	ID 				string 		`gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title 			string 		`gorm:"type:varchar(255);not null" json:"title"`
	Description 	string 		`gorm:"type:text" json:"description"`
	StartDate 		time.Time 	`gorm:"not null" json:"start_date"`
	EndDate 		time.Time 	`gorm:"not null" json:"end_date"`
	CreatorID 		string 		`gorm:"type:uuid;not null" json:"creator_id"`
	CreatedAt 		time.Time 	`gorm:"not null" json:"created_at"`
	UpdatedAt 		time.Time 	`gorm:"not null" json:"updated_at"`
}