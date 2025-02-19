package repository

import (
	"context"

	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/service"
)

type EventRepository interface {
    GetByID(ctx context.Context, id string) (*model.Event, error)
    List(ctx context.Context, limit, offset int, sortBy, sortDir string) ([]*model.Event, error)
    Create(ctx context.Context, event *model.Event) error
    Update(ctx context.Context, event *model.Event) error
    Delete(ctx context.Context, id string) error
    Search(ctx context.Context, params *service.SearchEventsInput) ([]*model.Event, int64, error)
}