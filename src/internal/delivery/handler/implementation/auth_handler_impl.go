package handlerImplementation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/hafiztri123/src/internal/delivery/handler"
	"github.com/hafiztri123/src/internal/pkg/response"
	"github.com/hafiztri123/src/internal/service"
)

type authHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

func NewAuthHandler(authService service.AuthService) handler.AuthHandler {
	return &authHandler{
		authService: authService,
		validator:   validator.New(),
	}
}

func (h *authHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input service.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "[FAIL] invalid request payload")
		return
	}

	if err := h.validator.Struct(input); err != nil {
		respondWithError(w, http.StatusBadRequest, "[FAIL] validation error")
	}


	err := h.authService.Register(&input)
	if err != nil {
		if err.Error() == "[FAIL] user already exists" {
			respondWithError(w, http.StatusConflict, fmt.Sprintf("[FAIL]  %s", err))
			return
		}
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("[FAIL] %s", err) )
	}

	respondWithJSON(w, http.StatusCreated, response.SuccessResponse{
		Timestamp: time.Now(),
		Data: "",
	})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithError(w, http.StatusBadRequest, "[FAIL] invalid request payload ")
		return
	}

	if err := h.validator.Struct(input); err != nil {
		respondWithError(w, http.StatusBadRequest, "[FAIL] validation error")
	}

	token, err := h.authService.Login(&input)
	if err != nil {
		if err.Error() == "[FAIL] user not found" || err.Error() == "[FAIL] invalid credentials" {
			respondWithError(w, http.StatusUnauthorized, "[FAIL] invalid credentials")
			return
		}

		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("[FAIL] %s", err))
	}

	respondWithJSON(w, http.StatusOK, response.SuccessResponse{
		Timestamp: time.Now(),
		Data: token,
	})
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	error := response.ErrorResponse{
		Timestamp: time.Now(),
		Message: message,
		

	}
	respondWithJSON(w, code, error)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
