package repository

import (
    "context"
    "errors"
    "fmt"
    "log"
    "time"

    "github.com/hafiztri123/src/internal/model"
    "github.com/hafiztri123/src/internal/pkg/cache"
    "gorm.io/gorm"
)

// EventRepository defines the interface for event-related database operations.
type EventRepository interface {
    GetByID(ctx context.Context, id string) (*model.Event, error)
    List(ctx context.Context, limit, offset int, sortBy, sortDir string) ([]*model.Event, error)
    Create(ctx context.Context, event *model.Event) error
    Update(ctx context.Context, event *model.Event) error
    Delete(ctx context.Context, id string) error
    Search(ctx context.Context, params *model.SearchEventsInput) ([]*model.Event, int64, error)
}

// eventRepository implements the EventRepository interface.
type eventRepository struct {
    db    *gorm.DB
    cache *cache.RedisCache
}

// NewEventRepository creates a new instance of EventRepository.
func NewEventRepository(db *gorm.DB, cache *cache.RedisCache) EventRepository {
    return &eventRepository{
        db:    db,
        cache: cache,
    }
}

// GetByID retrieves an event by its ID, using cache if available.
func (r *eventRepository) GetByID(ctx context.Context, id string) (*model.Event, error) {
    var event model.Event
    cacheKey := fmt.Sprintf("event:%s", id)
    err := r.cache.Get(ctx, cacheKey, &event)

    if err == nil && event.ID != "" {
        return &event, nil
    }

    result := r.db.Where("id = ?", id).First(&event)
    if result.Error != nil {
        if errors.Is(result.Error, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, result.Error
    }

    err = r.cache.Set(ctx, cacheKey, event, 30*time.Minute)
    if err != nil {
        log.Printf("[FAIL] failed to cache event: %v", err)
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

    result := query.Find(&events)
    if result.Error != nil {
        return nil, result.Error
    }

    err = r.cache.Set(ctx, cacheKey, events, 5*time.Minute)
    if err != nil {
        log.Printf("[FAIL] Failed to cache events list: %s", err)
    }

    return events, nil
}

// Create inserts a new event into the database and invalidates relevant cache keys.
func (r *eventRepository) Create(ctx context.Context, event *model.Event) error {
    err := r.db.Create(event).Error
    if err != nil {
        log.Printf("[FAIL] failed to cache new event: %v", err)
    }

    listKeysPattern := "events:list:*"
    keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
    if err != nil {
        log.Printf("[FAIL] Failed to get list cache keys: %v", err)
    }

    for _, key := range keys {
        err = r.cache.Delete(ctx, key)
        if err != nil {
            log.Printf("[FAIL] failed to invalidate list cache: %v", err)
        }
    }

    return nil
}

// Update updates an existing event in the database and invalidates relevant cache keys.
func (r *eventRepository) Update(ctx context.Context, event *model.Event) error {
    err := r.db.Transaction(func(tx *gorm.DB) error {
        err := tx.Save(event).Error
        if err != nil {
            return err
        }
        return nil
    })

    if err != nil {
        return err
    }

    cacheKey := fmt.Sprintf("events:%s", event.ID)
    err = r.cache.Set(ctx, cacheKey, event, 30*time.Minute)
    if err != nil {
        log.Printf("[FAIL] failed to update event cache: %v", err)
    }

    listKeysPattern := "events:list:*"
    keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
    if err != nil {
        log.Printf("[FAIL] failed to get list cache keys: %v", err)
    }

    for _, key := range keys {
        err = r.cache.Delete(ctx, key)
        if err != nil {
            log.Printf("[FAIL] failed to invalidate list cache: %v", err)
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
        return fmt.Errorf("[FAIL] failed to delete event: %s", err)
    }

    cacheKey := fmt.Sprintf("event:%s", id)
    err = r.cache.Delete(ctx, cacheKey)
    if err != nil {
        log.Printf("[FAIL] failed to delete event cache: %v", err)
    }

    listKeysPattern := "events:list:*"
    keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
    if err != nil {
        log.Printf("[FAIL] failed to get list cache keys: %v", err)
    }

    for _, key := range keys {
        err = r.cache.Delete(ctx, key)
        if err != nil {
            log.Printf("[FAIL] failed to invalidate list cache: %v", err)
        }
    }

    return nil
}

// Search retrieves events based on search parameters.
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
        return nil, 0, err
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
        return nil, 0, err
    }

    return events, totalCount, nil
}