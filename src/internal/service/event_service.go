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
	Page 		int `json:"page" validate:"min=1"`
	PageSize 	int `json:"page_size" validate:"min=1,max=100"`
}
