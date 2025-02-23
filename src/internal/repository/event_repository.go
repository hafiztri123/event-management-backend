package repository

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/pkg/cache"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"gorm.io/gorm"
)

type EventRepository interface {
    GetByID(ctx context.Context, id string) (*model.Event, error)
    List(ctx context.Context, limit, offset int, sortBy, sortDir string) ([]*model.Event, error)
    Create(ctx context.Context, event *model.Event) error
    Update(ctx context.Context, event *model.Event) error
    Delete(ctx context.Context, id string) error
    Search(ctx context.Context, params *model.SearchEventsInput) ([]*model.Event, int64, error)
}

type eventRepository struct {
    db    *gorm.DB
    cache *cache.RedisCache
}

func NewEventRepository(db *gorm.DB, cache *cache.RedisCache) EventRepository {
    return &eventRepository{
        db:    db,
        cache: cache,
    }
}

const (
    CACHE_SET_FAIL string = "[FAIL] Setting cached failed"
    CACHE_KEYS_FAIL string = "[FAIL] Getting cached keys failed"
    CACHE_DELETE_FAIL string = "[FAIL] Delete cached key failed"
)

func (r *eventRepository) GetByID(ctx context.Context, id string) (*model.Event, error) {
    var event model.Event
    cacheKey := fmt.Sprintf("event:%s", id)
    err := r.cache.Get(ctx, cacheKey, &event)

    if err == nil && event.ID != "" {
        return &event, nil
    }

    err = r.db.Where("id = ?", id).First(&event).Error
    if err != nil {
        return nil, errs.NewNotFoundError("Event not found")
    }

    err = r.cache.Set(ctx, cacheKey, event, 30*time.Minute)
    if err != nil {
        log.Printf("%s: %v", CACHE_SET_FAIL, err)
    }

    return &event, nil
}

// List retrieves a paginated list of events, using cache if available.
func (r *eventRepository) List(ctx context.Context, limit, offset int, sortBy, sortDir string) ([]*model.Event, error) {
    var events []*model.Event
    cacheKey := fmt.Sprintf("events:list:%d:%d:%s:%s", limit, offset, sortBy, sortDir)
    err := r.cache.Get(ctx, cacheKey, &events)

    if err == nil && len(events) > 0 {
        return events, nil
    }

    if sortBy == "" {
        sortBy = "created_at"
    }

    validSortColumns := map[string]bool{
        "title":      true,
        "start_date": true,
        "end_date":   true,
        "created_at": true,
    }

    if !validSortColumns[sortBy] {
        sortBy = "created_at"
    }

    if sortDir != "ASC" && sortDir != "asc" {
        sortDir = "DESC"
    } else {
        sortDir = "ASC"
    }

    query := r.db.
        Limit(limit).
        Offset(offset).
        Order(fmt.Sprintf("%s %s", sortBy, sortDir))

    err = query.Find(&events).Error
    if err != nil {
        return nil, DBError(err)
    }

    err = r.cache.Set(ctx, cacheKey, events, 5*time.Minute)
    if err != nil {
        log.Printf("%s: %s", CACHE_SET_FAIL, err)
    }

    return events, nil
}

// Create inserts a new event into the database and invalidates relevant cache keys.
func (r *eventRepository) Create(ctx context.Context, event *model.Event) error {
    err := r.db.Create(event).Error
    if err != nil {
        return DBError(err)
    }

    listKeysPattern := "events:list:*"
    keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
    if err != nil {
        log.Printf("%s: %v",CACHE_KEYS_FAIL, err)
    }

    for _, key := range keys {
        err = r.cache.Delete(ctx, key)
        if err != nil {
            log.Printf("%s: %v", CACHE_DELETE_FAIL, err)
        }
    }

    return nil
}

func (r *eventRepository) Update(ctx context.Context, event *model.Event) error {
    err := r.db.Transaction(func(tx *gorm.DB) error {
        err := tx.Save(event).Error
        if err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        return DBError(err)
    }

    cacheKey := fmt.Sprintf("events:%s", event.ID)
    err = r.cache.Set(ctx, cacheKey, event, 30*time.Minute)
    if err != nil {
        log.Printf("%s: %v", CACHE_SET_FAIL, err)
    }

    listKeysPattern := "events:list:*"
    keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
    if err != nil {
        log.Printf("%s: %v",CACHE_KEYS_FAIL, err)
    }

    for _, key := range keys {
        err = r.cache.Delete(ctx, key)
        if err != nil {
            log.Printf("%s: %v",CACHE_DELETE_FAIL, err)
        }
    }

    return nil
}

// Delete removes an event from the database and invalidates relevant cache keys.
func (r *eventRepository) Delete(ctx context.Context, id string) error {
    err := r.db.Transaction(func(tx *gorm.DB) error {
        err := tx.Delete(&model.Event{}, "id = ?", id).Error
        if err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        return DBError(err)
    }

    cacheKey := fmt.Sprintf("event:%s", id)
    err = r.cache.Delete(ctx, cacheKey)
    if err != nil {
        log.Printf("%s: %v",CACHE_DELETE_FAIL, err)
    }

    listKeysPattern := "events:list:*"
    keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
    if err != nil {
        log.Printf("%s: %v",CACHE_KEYS_FAIL, err)
    }

    for _, key := range keys {
        err = r.cache.Delete(ctx, key)
        if err != nil {
        log.Printf("%s: %v",CACHE_DELETE_FAIL, err)
        }
    }

    return nil
}

func (r *eventRepository) Search(ctx context.Context, params *model.SearchEventsInput) ([]*model.Event, int64, error) {
    var events []*model.Event
    var totalCount int64

    query := r.db.Model(&model.Event{})

    if params.Query != "" {
        query = query.Where("title ILIKE ? OR description ILIKE ?",
            "%"+params.Query+"%", "%"+params.Query+"%")
    }

    if params.StartDate != nil {
        query = query.Where("start_date >= ?", params.StartDate)
    }

    if params.EndDate != nil {
        query = query.Where("end_date <= ?", params.EndDate)
    }

    if params.Creator != "" {
        query = query.Where("creator_id = ?", params.Creator)
    }

    if err := query.Count(&totalCount).Error; err != nil {
        return nil, 0, DBError(err)
    }

    sortBy := "created_at"
    if params.SortBy != "" {
        switch params.SortBy {
        case "title", "start_date", "end_date", "created_at":
            sortBy = params.SortBy
        }
    }

    sortDir := "DESC"
    if params.SortDir == "asc" {
        sortDir = "ASC"
    }

    offset := (params.Page - 1) * params.PageSize
    err := query.Order(fmt.Sprintf("%s %s", sortBy, sortDir)).
        Limit(params.PageSize).
        Offset(offset).
        Find(&events).Error

    if err != nil {
        return nil, 0, DBError(err)
    }

    return events, totalCount, nil
}