package model

import "time"

type Event struct {
	ID 				string 		`json:"id"`
	Title 			string 		`json:"title"`
	Description 	string 		`json:"description"`
	StartDate 		time.Time 	`json:"start_date"`
	EndDate 		time.Time 	`json:"end_date"`
	CreatorID 		string 		`json:"creator_id"`
	CreatedAt 		time.Time 	`json:"created_at"`
	UpdatedAt 		time.Time 	`json:"updated_at"`
}