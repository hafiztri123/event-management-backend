package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/pkg/response"
	"github.com/hafiztri123/src/internal/service"
)


type UserHandler interface {
    UpdateProfile(w http.ResponseWriter, r *http.Request)
    GetProfile(w http.ResponseWriter, r *http.Request)
    ChangePassword(w http.ResponseWriter, r *http.Request)
    // UploadProfileImage(w http.ResponseWriter, r *http.Request)
}

type userHandlerImpl struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) UserHandler {
	return userHandlerImpl{
		userService: userService,
	}
}

func (h userHandlerImpl) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var input model.UpdateProfileInput
	err := json.NewDecoder(r.Body).Decode(&input)
	token := r.Context().Value("user").(jwt.MapClaims)
	userID := token["user_id"].(string)


	if err != nil {
		respondWithJSON(w, 400, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] fail to parse request. Bad request",
		})
	}

	err = h.userService.UpdateProfile(userID, &input)
	if err != nil {
		respondWithJSON(w, 500, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] fail to update profile. Bad request",
		})
	}

	respondWithJSON(w, 200, response.Response{
		Timestamp: time.Now(),
	})
}


func (h userHandlerImpl) GetProfile(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("user").(jwt.MapClaims)["user_id"].(string)
	if userID == "" {
		respondWithJSON(w, 401, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Unauthorized",
		})
	}

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
		respondWithJSON(w, 500, response.Response{
			Timestamp: time.Now(),
			Message: fmt.Sprintf("[FAIL] fail to get profile: %v", err),
		})
	}

	respondWithJSON(w, 200, response.Response{
		Timestamp: time.Now(),
		Data: profile,
	})
 
}

func (h userHandlerImpl) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user").(jwt.MapClaims)["user_id"].(string)
	if userID == "" {
		respondWithJSON(w, 401, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Unauthorized",
		})
	}

	var input model.ChangePasswordInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		respondWithJSON(w, 400, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] fail to parse request. Bad request",
		})
	}

	err = h.userService.ChangePassword(userID, &input)
	if err != nil {
		respondWithJSON(w, 500, response.Response{
			Timestamp: time.Now(),
			Message: fmt.Sprintf("[FAIL] fail to change password: %v", err),
		})
	}

	respondWithJSON(w, 200, response.Response{
		Timestamp: time.Now(),
	})
}