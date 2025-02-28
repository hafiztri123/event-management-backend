package model

import "time"

type Event struct {
	ID 				string 		`gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Title 			string 		`gorm:"type:varchar(255);not null" json:"title"`
	Description 	string 		`gorm:"type:text" json:"description"`
	StartDate 		time.Time 	`gorm:"not null" json:"start_date"`
	EndDate 		time.Time 	`gorm:"not null" json:"end_date"`
	CreatorID 		string 		`gorm:"type:uuid;not null" json:"creator_id"`
	CategoryID 		string 		`gorm:"type:uuid" json:"category_id"`
	Tags 			[]Tag		`gorm:"many2many:event_tags" json:"tags"`
	Files 			[]File		`gorm:"foreignKey:EventID" json:"files"`
	CreatedAt 		time.Time 	`gorm:"not null" json:"created_at"`
	UpdatedAt 		time.Time 	`gorm:"not null" json:"updated_at"`
}

type Tag struct {
	ID 		string 	`gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name	string 	`gorm:"type:varchar(50);not null;unique" json:"name"`
	Events 	[]Event `gorm:"many2many:event_tags;" json:"-"`
}

type File struct {
	ID 			string 		`gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	EventID 	string 		`gorm:"type:uuid;not null" json:"event_id"`
	FileName 	string 		`gorm:"type:varchar(255);not null" json:"file_name"`
	FileType 	string 		`gorm:"type:varchar(100);not null" json:"file_type"`
	FileURL 	string 		`gorm:"type:text;not null" json:"file_url"`
	CreatedAt 	time.Time 	`gorm:"not null" json:"created_at"`
}

type CreateEventInput struct {
	Title 		string 		`json:"title" validate:"required"`
	Description string 		`json:"description"`
	CategoryID 	string 		`json:"category_id"`
	StartDate 	time.Time 	`json:"start_date" validate:"required"`
	EndDate 	time.Time 	`json:"end_date" validate:"required,gtfield=StartDate"`
}

type UpdateEventInput struct {
	Title 		string 		`json:"title" validate:"required"`
	Description string 		`json:"description"`
	StartDate 	time.Time 	`json:"start_date" validate:"required"`
	EndDate 	time.Time 	`json:"end_date" validate:"required,gtfield=StartDate"`
}

type ListEventsInput struct {
	Page 		int 	`json:"page" validate:"min=1"`
	PageSize 	int 	`json:"page_size" validate:"min=1,max=100"`
	SortBy 		string 	`json:"sort_by,omitempty"`
	SortDir 	string 	`json:"sort_dir,omitempty"`
}

type SearchEventsInput struct {
	Query 		string 		`json:"query,omitempty"`
	StartDate 	*time.Time 	`json:"start_date,omitempty"`
	EndDate 	*time.Time 	`json:"end_date,omitempty"`
	Creator 	string 		`json:"creator,omitempty"`
	Page 		int 		`json:"page" validate:"min=1"`
	PageSize 	int 		`json:"page_size" validate:"min=1,max=100"`
	SortBy 		string 		`json:"sort_by,omitempty"`
	SortDir 	string 		`json:"sort_dir,omitempty"`
}

type SearchEventsOutput struct {
	Events 		[]*Event 		`json:"events"`
	TotalCount 	int64 			`json:"total_count"`
	Page 		int 			`json:"page"`
	PageSize 	int 			`json:"page_size"`
	TotalPages 	int 			`json:"total_pages"`
}

type UploadFile struct {
	FileName 	string 		`gorm:"type:varchar(255);not null" json:"file_name"`
	FileType 	string 		`gorm:"type:varchar(100);not null" json:"file_type"`
}


	