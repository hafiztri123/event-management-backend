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


type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	Update(ctx context.Context, category *model.Category) error
	Delete(ctx context.Context, id string) error
	GetByID(ctx context.Context, id string) (*model.Category, error)
	List(ctx context.Context) ([]*model.Category, error)
	IsIDExists(id string) (bool)
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
		return DBError(err)
	}

	listKeysPattern := "categories:list:*"
	keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
	if err != nil {
		log.Printf("%s: %v", CACHE_KEYS_FAIL, err)
	}

	for _, key := range keys {
		err = r.cache.Delete(ctx, key)
		if err != nil {
			log.Printf("%s: %v", CACHE_DELETE_FAIL, err)
		} 
	}

	return nil
}

func (r *categoryRepositoryImpl) Update(ctx context.Context, category *model.Category) error {
	if category == nil {
		return errs.NewBadRequestError("")
	}

	var existingCategory model.Category
	model := r.db.Model(&model.Category{})
	err := model.Where("id = ?", category.ID).First(&existingCategory).Error
	if err != nil {
		return DBError(err)
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
		return DBError(err)
	}

	cacheKey := fmt.Sprintf("categories:%s", existingCategory.ID)
	err = r.cache.Set(ctx, cacheKey, existingCategory, 30*time.Minute)
	if err != nil {
		log.Printf("%s: %v", CACHE_SET_FAIL, err)
	}

	listKeysPattern := "events:list:*"
	keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
	if err != nil {
		log.Printf("%s: %v", CACHE_KEYS_FAIL ,err)
	}

	for _, key := range keys {
		err := r.cache.Delete(ctx, key)
		if err != nil {
			log.Printf("%s: %v",CACHE_DELETE_FAIL, err)
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
		return DBError(err)
	}

	cacheKey := fmt.Sprintf("categories:%s", categoryID)
	err = r.cache.Delete(ctx, cacheKey)
	if err != nil {
		log.Printf("%s: %v",CACHE_DELETE_FAIL,err)
	}

	listKeysPattern := "categories:list:*"
	keys, err := r.cache.Client.Keys(ctx, listKeysPattern).Result()
	if err != nil {
		log.Printf("%s: %v",CACHE_KEYS_FAIL, err)
	}

	for _, key := range keys {
		err = r.cache.Delete(ctx, key)
		if err != nil {
			return fmt.Errorf("%s: %v",CACHE_DELETE_FAIL, err)
		}
	}

	return nil

}
 
func (r *categoryRepositoryImpl) GetByID(ctx context.Context, id string) (*model.Category, error){
	var existingCategory model.Category
	category := r.db.Model(&model.Category{})
	err := category.Where("id = ?", id).First(&existingCategory).Error
	if err != nil {
		return nil,  DBError(err)
	}

	cacheKey := fmt.Sprintf("categories:%s", id)
	err = r.cache.Set(ctx, cacheKey, existingCategory, 30*time.Minute)

	if err != nil {
		log.Printf("%s: %v",CACHE_SET_FAIL, err)
	}
	return &existingCategory, nil
}

func (r *categoryRepositoryImpl) List(ctx context.Context) ([]*model.Category, error){
	var existingCategories []*model.Category
	category := r.db.Model(&model.Category{})
	err := category.Find(&existingCategories).Error
	if err != nil{
		return nil, DBError(err)
	}

	listKeysPattern := "categories:list"
	err = r.cache.Set(ctx, listKeysPattern, existingCategories, 30*time.Minute)

	if err != nil {
		log.Printf("%s: %v",CACHE_SET_FAIL, err)
	}

	return existingCategories, nil
} 

func (r *categoryRepositoryImpl) IsIDExists(id string) (bool) {
	var idCount int64
	model := r.db.Model(&model.Category{})
	model.Where("id = ?", id).Count(&idCount)
	return idCount > 0
}







