package repositoryImplementation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/pkg/cache"
	"github.com/hafiztri123/src/internal/repository"
	"gorm.io/gorm"
)

type eventRepository struct {
	db *gorm.DB
	cache *cache.RedisCache
}

func NewEventRepository(db *gorm.DB, cache *cache.RedisCache) repository.EventRepository {
	return &eventRepository{
		db: db,
		cache: cache,

	}
}

func (r *eventRepository) GetByID(ctx context.Context, id string) (*model.Event, error){
	var event model.Event

	cacheKey := fmt.Sprintf("event:%s", id)
	err := r.cache.Get(ctx, cacheKey, &event)
	
	if err == nil && event.ID != "" {
		return &event, nil
	}

	result := r.db.Where("id = ?", id).First(&event)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound){
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

func (r *eventRepository) List(ctx context.Context, limit, offset int) ([]*model.Event, error) {
	var events []*model.Event

	cacheKey := fmt.Sprintf("events:list:%d:%d", limit, offset)
	err := r.cache.Get(ctx, cacheKey, &events)
	if err == nil && len(events) > 0 {
		return events, nil
	}


	result := r.db.
		Limit(limit).
		Offset(offset).
		Order("created_at DESC").
		Find(&events)
	
	if result.Error != nil {
		return nil, result.Error
	}

	err = r.cache.Set(ctx, cacheKey, events, 5*time.Minute)
	if err != nil {
		log.Printf("[FAIL] Failed to cache events list: %s", err)
	}

	return events, nil
}

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

func (r *eventRepository) Delete(ctx context.Context, id string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&model.Event{}, "id = ?", id).Error
		if err != nil{
			return err
		}
		return nil
	})

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


