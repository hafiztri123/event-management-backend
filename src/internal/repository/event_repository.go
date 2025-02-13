package repository

import "github.com/hafiztri123/src/internal/model"

type EventRepository interface {
	GetByID(id string) (*model.Event, error)
    List(limit, offset int) ([]*model.Event, error)
    Create(event *model.Event) error
    Update(event *model.Event) error
    Delete(id string) error
}