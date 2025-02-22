package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/pkg/cache"
	"gorm.io/gorm"
)


type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Category, error)
	List(ctx context.Context) ([]*model.Category, error)
	IsIDExists(id string) (bool, error)
}

type categoryRepositoryImpl struct {
	db *gorm.DB
	cache *cache.RedisCache
}

func NewCategoryRepository(db *gorm.DB,  redis *cache.RedisCache) CategoryRepository {
	return &categoryRepositoryImpl{
		db: db,
		cache: redis,

	}
}

func(r *categoryRepositoryImpl) Create(ctx context.Context, category *model.Category) error {
	model := r.db.Model(&model.Category{})
	
	err := model.Create(category).Error
	if err != nil {
		return fmt.Errorf("[FAIL] Failed to create event: %v", err)
	}

	listKeysPattern := "categories:list:*"
	keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
	if err != nil {
		return fmt.Errorf("[FAIL] Failed to get list cache keys: %v", err)
	}

	for _, key := range keys {
		err = r.cache.Delete(ctx, key)
		if err != nil {
			return fmt.Errorf("[FAIL] failed to invalidate list cache: %v", err)
		} 
	}

	return nil
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, category *model.Category) error {
	if category == nil {
		return fmt.Errorf("[FAIL] Update request is empty")
	}

	var existingCategory model.Category
	model := r.db.Model(&model.Category{})
	err := model.Where("id = ?", category.ID).First(&existingCategory).Error
	if err != nil {
		return fmt.Errorf("[FAIL] failed to get existing category: %v", err)
	}

	if category.Name != "" {
		existingCategory.Name = category.Name
	}

	if category.Description != "" {
		existingCategory.Description = category.Description
	}

	if category.Name != "" || category.Description != "" {
		existingCategory.UpdatedAt = time.Now()
	} 

	err = r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Save(existingCategory).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("[FAIL] failed to update existing category: %v", err)
	}

	cacheKey := fmt.Sprintf("categories:%s", existingCategory.ID)
	err = r.cache.Set(ctx, cacheKey, existingCategory, 30*time.Minute)
	if err != nil {
		return fmt.Errorf("[FAIL] fail to update categories cache: %v", err)
	}

	listKeysPattern := "events:list:*"
	keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
	if err != nil {
		return fmt.Errorf("[FAIL] failed to get list cache categories keys: %v", err)
	}

	for _, key := range keys {
		err := r.cache.Delete(ctx, key)
		if err != nil {
			return fmt.Errorf("[FAIL] failed to invalidate cache categories keys: %v", err)
		}
	}

	return nil
}

func (r *categoryRepositoryImpl) Delete(ctx context.Context, categoryID string) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Delete(&model.Category{}, "id = ?", categoryID).Error
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("[FAIL] fail to delete transaction: %v", err)
	}

	cacheKey := fmt.Sprintf("categories:%s", categoryID)
	err = r.cache.Delete(ctx, cacheKey)
	if err != nil {
		return fmt.Errorf("[FAIL] failed to delete categories cache: %v", err)
	}

	listKeysPattern := "categories:list:*"
	keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
	if err != nil {
		return fmt.Errorf("[FAIL] failed to get list cache keys: %v", err)
	}

	for _, key := range keys {
		err = r.cache.Delete(ctx, key)
		if err != nil {
			return fmt.Errorf("[FAIL] failed to invalidate list cache: %v", err)
		}
	}

	return nil

}
 
func (r *categoryRepositoryImpl) GetByID(ctx context.Context, id string) (*model.Category, error){
	var existingCategory model.Category
	category := r.db.Model(&model.Category{})
	err := category.Where("id = ?", id).First(&existingCategory).Error
	if err != nil {
		return nil,  fmt.Errorf("[FAIL] failed to get category: %v", err)
	}

	cacheKey := fmt.Sprintf("categories:%s", id)
	err = r.cache.Set(ctx, cacheKey, existingCategory, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("[FAIL] failed to set cache for categories: %v", err)
	}
	return &existingCategory, nil
}

func (r *categoryRepositoryImpl) List(ctx context.Context) ([]*model.Category, error){
	var existingCategories []*model.Category
	category := r.db.Model(&model.Category{})
	err := category.Find(&existingCategories).Error
	if err != nil{
		return nil, fmt.Errorf("[FAIL] failed to get list of categories: %v", err)
	}

	listKeysPattern := "categories:list"
	err = r.cache.Set(ctx, listKeysPattern, existingCategories, 30*time.Minute)
	if err != nil {
		return nil, fmt.Errorf("[FAIL] failed to set list of categories cache: %v", err)
	}

	return existingCategories, nil
} 

func (r *categoryRepositoryImpl) IsIDExists(id string) (bool, error) {
	var idCount int64
	model := r.db.Model(&model.Category{})
	err := model.Where("id = ?", id).Count(&idCount)
	if err != nil {
		return false, nil
	}
	return idCount > 0, nil
}






