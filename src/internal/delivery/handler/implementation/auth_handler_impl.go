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
		respondWithJSON(w, http.StatusBadRequest, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Attempt to parse request has failed. Bad request.",
		})
		return
	}

	if err := h.validator.Struct(input); err != nil {
		respondWithJSON(w, http.StatusBadRequest, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Request was not valid. Bad request.",
		})
	}


	err := h.authService.Register(&input)
	if err != nil {
		if err.Error() == "[FAIL] user already exists" {
			respondWithJSON(w, http.StatusConflict, response.Response{
				Timestamp: time.Now(),
				Message: "[FAIL] User already exists. Conflict.",
			})
			return
		}
		respondWithJSON(w, http.StatusInternalServerError, response.Response{
			Timestamp: time.Now(),
			Message: fmt.Sprintf("[FAIL] %s. Internal server error.", err),
		} )
	}

	respondWithJSON(w, http.StatusCreated, response.Response{
		Timestamp: time.Now(),
	})
}

func (h *authHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input service.LoginInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithJSON(w, http.StatusBadRequest, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Attempt to parse request has failed. Bad request.",
		})
		return
	}

	if err := h.validator.Struct(input); err != nil {
		respondWithJSON(w, http.StatusBadRequest, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Request was not valid. Bad request",
		})
	}

	loginResponse, err := h.authService.Login(&input)
	if err != nil {
		if err.Error() == "[FAIL] user not found" || err.Error() == "[FAIL] invalid credentials" {
			respondWithJSON(w, http.StatusUnauthorized, response.Response{
				Timestamp: time.Now(),
				Message: "[FAIL] Invalid credentials. Unauthorized",
			})
			return
		}

		respondWithJSON(w, http.StatusInternalServerError, response.Response{
			Timestamp: time.Now(),
			Message: fmt.Sprintf("[FAIL] %s. Internal server error", err),
		})
	}

	respondWithJSON(w, http.StatusOK, response.Response{
		Timestamp: time.Now(),
		Data: loginResponse,
	})

}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
