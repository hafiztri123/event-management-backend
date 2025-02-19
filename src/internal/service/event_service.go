package service

import (
	"time"

	"github.com/hafiztri123/src/internal/model"
)


type EventService interface {
	CreateEvent(input *CreateEventInput, creatorID string) error 
	UpdateEvent(id string, input *UpdateEventInput, userID string) error
	DeleteEvent(id string, userID string) error
	GetEvent(id string) (*model.Event, error)
	ListEvents(input *ListEventsInput) ([]*model.Event, error)
	SearchEvents(input *SearchEventsInput) (*SearchEventsOutput, error)
}



type CreateEventInput struct {
	Title 		string 		`json:"title" validate:"required"`
	Description string 		`json:"description"`
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
	Events 		[]*model.Event 	`json:"events"`
	TotalCount 	int64 			`json:"total_count"`
	Page 		int 			`json:"page"`
	PageSize 	int 			`json:"page_size"`
	TotalPages 	int 			`json:"total_pages"`
}

	
