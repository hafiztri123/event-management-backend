package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"github.com/hafiztri123/src/internal/pkg/response"
	"github.com/hafiztri123/src/internal/service"
)


type CategoryHandler interface {
    CreateCategory(w http.ResponseWriter, r *http.Request)
    UpdateCategory(w http.ResponseWriter, r *http.Request)
    DeleteCategory(w http.ResponseWriter, r *http.Request)
    GetCategory(w http.ResponseWriter, r *http.Request)
    ListCategories(w http.ResponseWriter, r *http.Request)
}

type categoryHandlerImpl struct {
	categoryService service.CategoryService
	validator *validator.Validate
}

func NewCategoryHandler(categoryService service.CategoryService) CategoryHandler {
	return &categoryHandlerImpl{
		categoryService: categoryService,
		validator: validator.New(),
	}
}

func (h *categoryHandlerImpl) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var input model.CreateCategoryInput


	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        HandleErrorResponse(w, err)
		return
	}

	if err := h.validator.Struct(&input); err != nil {
		HandleErrorResponse(w, errs.NewValidationError("Request is not valid"))
		return 
	}


	err := h.categoryService.CreateCategory(&input)
	if err != nil {
        HandleErrorResponse(w, err)
		return
	}

	respondWithJSON(w, 201, response.Response{
		Timestamp: time.Now(),
	})
}

func (h *categoryHandlerImpl)  UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	var input model.UpdateCategoryInput

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
        HandleErrorResponse(w, err)
		return
	}

	err = h.categoryService.UpdateCategory(categoryID, &input)
	if err != nil {
        HandleErrorResponse(w, err)
		return
	}

	respondWithJSON(w, 201, response.Response{
		Timestamp: time.Now(),
	})


}
func (h *categoryHandlerImpl)  DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	
	err := h.categoryService.DeleteCategory(categoryID)
	if err != nil {
        HandleErrorResponse(w, err)
		return
	}

	respondWithJSON(w, 204, response.Response{
		Timestamp: time.Now(),
	})

}
func (h *categoryHandlerImpl)  GetCategory(w http.ResponseWriter, r *http.Request) {
	categoryID := chi.URLParam(r, "id")
	
	category, err := h.categoryService.GetCategory(categoryID)
	if err != nil {
        HandleErrorResponse(w, err)
		return
	}

	respondWithJSON(w, 200, response.Response{
		Timestamp: time.Now(),
		Data: category,
	})

}
func (h *categoryHandlerImpl)  ListCategories(w http.ResponseWriter, r *http.Request){
	categories, err := h.categoryService.ListCategories()
	if err != nil {
        HandleErrorResponse(w, err)
		return
	}

	respondWithJSON(w, 200, response.Response{
		Timestamp: time.Now(),
		Data: categories,
	})

}