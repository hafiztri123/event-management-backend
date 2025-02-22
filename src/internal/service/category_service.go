package service

import (
	"context"
	"fmt"
	"time"

	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/repository"
)

type CategoryService interface {
    CreateCategory(input *model.CreateCategoryInput) error
    UpdateCategory(id string, input *model.UpdateCategoryInput) error
    DeleteCategory(id string) error
    GetCategory(id string) (*model.Category, error)
    ListCategories() ([]*model.Category, error)
}



type categoryServiceImpl struct {
	categoryRepository repository.CategoryRepository
}

func NewCategoryService(categoryRepository repository.CategoryRepository) CategoryService {
	return categoryServiceImpl{
		categoryRepository: categoryRepository,
	}
}

func (s categoryServiceImpl) CreateCategory(input *model.CreateCategoryInput) error {
	if input == nil {
		return fmt.Errorf("[FAIL] input is missing")
	}

	ctx := context.Background()
	category := &model.Category{
		Name: input.Name,
		Description: input.Description,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	err := s.categoryRepository.Create(ctx, category)
	if err != nil {
		return fmt.Errorf("[FAIL] failed creating category: %s", err)
	}

	return nil
}

func (s categoryServiceImpl) UpdateCategory(id string, input *model.UpdateCategoryInput) error {
	if input == nil {
		return fmt.Errorf("[FAIL] input is missing")
	}

	exists, err := s.isCategoryIDExists(id)
	if err != nil {
		return fmt.Errorf("[FAIL] category id check invalid")
	}

	if !exists {
		return fmt.Errorf("[FAIL] id is invalid")
	}


	ctx := context.Background()

	updatedModel := &model.Category{
		ID: id,
		Name: input.Name,
		Description: input.Description,
		CreatedAt: time.Now(), //placeholder
		UpdatedAt: time.Now(), //placeholder
	}
	
	if err := s.categoryRepository.Update(ctx, updatedModel); err != nil {
		return fmt.Errorf("[FAIL] failed to update category: %v", err)
	}

	return nil

}

func (s categoryServiceImpl) DeleteCategory(id string) error {
	if id == "" {
		return fmt.Errorf("[FAIL] id is missing")
	}

	exists, err := s.isCategoryIDExists(id)
	if err != nil {
		return fmt.Errorf("[FAIL] category id check invalid")
	}

	if !exists {
		return fmt.Errorf("[FAIL] id is invalid")
	}


	ctx := context.Background()

	if err := s.categoryRepository.Delete(ctx, id); err != nil {
		return fmt.Errorf("[FAIL] fail to delete category: %v", err)
	}

	return nil
}

func (s categoryServiceImpl) GetCategory(id string) (*model.Category, error){
	if id == "" {
		return nil, fmt.Errorf("[FAIL] id is missing")
	}

	exists, err := s.isCategoryIDExists(id)
	if err != nil {
		return nil, fmt.Errorf("[FAIL] category id check invalid")
	}

	if !exists {
		return nil, fmt.Errorf("[FAIL] id is invalid")
	}


	ctx := context.Background()

	category, err := s.categoryRepository.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("[FAIL] failed to get category: %v", err)
	}

	return category, nil
}


func (s categoryServiceImpl) ListCategories() ([]*model.Category, error){
	ctx := context.Background()

	categories, err := s.categoryRepository.List(ctx)
	
	if err != nil {
		return nil, fmt.Errorf("[FAIL] failed to get catogories: %v", err)
	}

	return categories, nil
}
 

func (s categoryServiceImpl) isCategoryIDExists(id string) (bool, error){
	return s.categoryRepository.IsIDExists(id)
}