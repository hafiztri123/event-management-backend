package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"github.com/hafiztri123/src/internal/pkg/response"
	"github.com/hafiztri123/src/internal/service"
)


type UserHandler interface {
    UpdateProfile(w http.ResponseWriter, r *http.Request)
    GetProfile(w http.ResponseWriter, r *http.Request)
    ChangePassword(w http.ResponseWriter, r *http.Request)
    UploadProfileImage(w http.ResponseWriter, r *http.Request)
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
        HandleErrorResponse(w, err)
        return
	}

	err = h.userService.UpdateProfile(userID, &input)
	if err != nil {
        HandleErrorResponse(w, err)
        return
	}

	respondWithJSON(w, 200, response.Response{
		Timestamp: time.Now(),
	})
}


func (h userHandlerImpl) GetProfile(w http.ResponseWriter, r *http.Request){
	userID := r.Context().Value("user").(jwt.MapClaims)["user_id"].(string)
	if userID == "" {
		HandleErrorResponse(w, errs.NewUnauthorizedError("Unauthorized"))
		return
	}

	profile, err := h.userService.GetProfile(userID)
	if err != nil {
        HandleErrorResponse(w, err)
        return
	}

	respondWithJSON(w, 200, response.Response{
		Timestamp: time.Now(),
		Data: profile,
	})
 
}

func (h userHandlerImpl) ChangePassword(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user").(jwt.MapClaims)["user_id"].(string)
	if userID == "" {
		HandleErrorResponse(w, errs.NewUnauthorizedError("Unauthorized"))
		return
	}

	var input model.ChangePasswordInput
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
        HandleErrorResponse(w, err)
        return
	}

	err = h.userService.ChangePassword(userID, &input)
	if err != nil {
        HandleErrorResponse(w, err)
        return
	}

	respondWithJSON(w, 200, response.Response{
		Timestamp: time.Now(),
	})
}

func (h userHandlerImpl) UploadProfileImage(w http.ResponseWriter, r *http.Request) {
	UserID := r.Context().Value("user").(jwt.MapClaims)["user_id"].(string)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		HandleErrorResponse(w, errs.NewBadRequestError("File too large"))
		return
	}

	

	file, header, err := r.FormFile("image")
	if err != nil {
		HandleErrorResponse(w, errs.NewBadRequestError("Invalid file"))
		return
	}
	defer file.Close()

	if header.Size > int64(10 << 20) { 
        HandleErrorResponse(w, errs.NewEntityTooLargeError("File too large"))
        return
    }

	if !isValidImageType(header.Header.Get("Content-Type")) {
		HandleErrorResponse(w, errs.NewBadRequestError("Invalid file type"))
		return
	}

	if err := h.userService.UploadProfileImage(r.Context(), UserID, file, header.Filename); err != nil {
		HandleErrorResponse(w, err)
	}

	respondWithJSON(w, http.StatusOK, response.Response{
		Timestamp: time.Now(),
	})

}

func isValidImageType(contentType string) bool {
    return contentType == "image/jpeg" || 
           contentType == "image/png" || 
           contentType == "image/gif"
}
