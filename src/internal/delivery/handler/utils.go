package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	errs "github.com/hafiztri123/src/internal/pkg/error"
)

// Response defines the structure of the API response.
type Response struct {
    Timestamp time.Time `json:"timestamp"`
    Code      int       `json:"code"`
    Message   string    `json:"message"`
}

// HandleErrorResponse handles errors and sends a standardized JSON response.
func HandleErrorResponse(w http.ResponseWriter, err error) {
    var badRequestErr *errs.BadRequestError
    if errors.As(err, &badRequestErr) {
        respondWithJSON(w, badRequestErr.Code, Response{
            Timestamp: time.Now(),
            Code:      badRequestErr.Code,
            Message:   badRequestErr.Message,
        })
        return
    }

    var notFoundErr *errs.NotFoundError
    if errors.As(err, &notFoundErr) {
        respondWithJSON(w, notFoundErr.Code, Response{
            Timestamp: time.Now(),
            Code:      notFoundErr.Code,
            Message:   notFoundErr.Message,
        })
        return
    }

    var validationErr *errs.ValidationError
    if errors.As(err, &validationErr) {
        respondWithJSON(w, validationErr.Code, Response{
            Timestamp: time.Now(),
            Code:      validationErr.Code,
            Message:   validationErr.Message,
        })
        return
    }

    var duplicateErr *errs.DuplicateEntryError
    if errors.As(err, &duplicateErr) {
        respondWithJSON(w, duplicateErr.Code, Response{
            Timestamp: time.Now(),
            Code:      duplicateErr.Code,
            Message:   duplicateErr.Message,
        })
        return
    }

    var forbiddenErr *errs.ForbiddenError
    if errors.As(err, &forbiddenErr) {
        respondWithJSON(w, forbiddenErr.Code, Response{
            Timestamp: time.Now(),
            Code:      forbiddenErr.Code,
            Message:   forbiddenErr.Message,
        })
        return
    }

    var databaseErr *errs.DatabaseError
    if errors.As(err, &databaseErr) {
        respondWithJSON(w, databaseErr.Code, Response{
            Timestamp: time.Now(),
            Code:      databaseErr.Code,
            Message:   databaseErr.Message,
        })
        return
    }

    var unauthorizedErr  *errs.UnauthorizedError
    if errors.As(err, &unauthorizedErr) {
        respondWithJSON(w, unauthorizedErr.Code, Response{
            Timestamp: time.Now(),
            Code: unauthorizedErr.Code,
            Message: unauthorizedErr.Message,
        })
        return
    }

    var EntityTooLargeErr  *errs.EntityTooLargeError
    if errors.As(err, &EntityTooLargeErr) {
        respondWithJSON(w, unauthorizedErr.Code, Response{
            Timestamp: time.Now(),
            Code: unauthorizedErr.Code,
            Message: unauthorizedErr.Message,
        })
        return
    }


    // Default response for unexpected errors
    respondWithJSON(w, http.StatusInternalServerError, Response{
        Timestamp: time.Now(),
        Code:      http.StatusInternalServerError,
        Message:   "An unexpected error occurred",
    })
}

func respondWithJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(statusCode)
    json.NewEncoder(w).Encode(payload)
}