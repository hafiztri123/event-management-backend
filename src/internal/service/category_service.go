package service

import (
	"context"
	"time"

	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
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
		return errs.NewBadRequestError("Request is missing")
	}

	ctx := context.Background()
	category := &model.Category{
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := s.categoryRepository.Create(ctx, category)
	if err != nil {
		return err
	}

	return nil
}

func (s categoryServiceImpl) UpdateCategory(id string, input *model.UpdateCategoryInput) error {
	if input.Description == "" && input.Name =="" {
		return errs.NewBadRequestError("Request is missing")
	}

	exists := s.isCategoryIDExists(id)

	if !exists {
		return errs.NewNotFoundError("Category not found")
	}

	ctx := context.Background()

	updatedModel := &model.Category{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		CreatedAt:   time.Now(), //placeholder
		UpdatedAt:   time.Now(), //placeholder
	}

	if err := s.categoryRepository.Update(ctx, updatedModel); err != nil {
		return err
	}

	return nil

}

func (s categoryServiceImpl) DeleteCategory(id string) error {
	if id == "" {
		return errs.NewBadRequestError("Category ID is missing")
	}

	exists := s.isCategoryIDExists(id)


	if !exists {
		return errs.NewNotFoundError("Category not found")
	}

	ctx := context.Background()

	if err := s.categoryRepository.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

func (s categoryServiceImpl) GetCategory(id string) (*model.Category, error) {
	if id == "" {
		return nil, errs.NewBadRequestError("Category ID is missing")
	}

	exists := s.isCategoryIDExists(id)


	if !exists {
		return nil, errs.NewNotFoundError("Category not found")
	}

	ctx := context.Background()

	category, err := s.categoryRepository.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return category, nil
}

func (s categoryServiceImpl) ListCategories() ([]*model.Category, error) {
	ctx := context.Background()

	categories, err := s.categoryRepository.List(ctx)

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (s categoryServiceImpl) isCategoryIDExists(id string) bool {
	return s.categoryRepository.IsIDExists(id)
}
