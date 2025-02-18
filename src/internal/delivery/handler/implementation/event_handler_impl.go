package handlerImplementation

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiztri123/src/internal/delivery/handler"
	"github.com/hafiztri123/src/internal/pkg/response"
	"github.com/hafiztri123/src/internal/service"
)

type eventHandler struct {
	eventService service.EventService
	validator *validator.Validate
}

func NewEventHandler(eventService service.EventService) handler.EventHandler {
	return &eventHandler{
		eventService: eventService,
		validator: validator.New(),
	}
}


func (h *eventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	var input service.CreateEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithJSON(w, http.StatusBadRequest, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] attempt to parse request has failed. Bad request.",
		})
		return
	}

	if err := h.validator.Struct(input); err != nil {
		respondWithJSON(w, http.StatusBadRequest, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Request was not valid. Bad request.",
		})
	}

	claims := r.Context().Value("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	err := h.eventService.CreateEvent(&input, userID)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, response.Response{
			Timestamp: time.Now(),
			Message: fmt.Sprintf("[FAIL] %s. Status internal server error", err.Error()),
		})
	}

	respondWithJSON(w, http.StatusCreated, response.Response{
		Timestamp: time.Now(),
	})

}

func (h *eventHandler) UpdateEvent (w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	var input service.UpdateEventInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		respondWithJSON(w, http.StatusBadRequest, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Attempt to parse request has failed. Bad request ",
		})
	}

	if err := h.validator.Struct(input); err != nil {
		respondWithJSON(w, http.StatusBadRequest, response.Response{
			Timestamp: time.Now(),
			Message: "[FAIL] Request was not valid. Bad request.",
		})
	}

	claims := r.Context().Value("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	err := h.eventService.UpdateEvent(eventID, &input, userID)
	if err != nil {
		switch err.Error() {
		case "[FAIL] event not found":
			respondWithJSON(w, http.StatusNotFound, fmt.Sprintf("%s. Not found.", response.Response{
				Timestamp: time.Now(),
				Message: fmt.Sprintf("%s. Not found.", err.Error()),
			}))
		case "[FAIL] unauthorized to modify this event":
			respondWithJSON(w, http.StatusForbidden, fmt.Sprintf("%s. Forbidden.", response.Response{
				Timestamp: time.Now(),
				Message: fmt.Sprintf("%s. Forbidden.", err.Error()),
			}))
		default:
			respondWithJSON(w, http.StatusInternalServerError, response.Response{
				Timestamp: time.Now(),
				Message: fmt.Sprintf("[FAIL] %s. Internal server error", err.Error()),
			})
		}

	}
	
	respondWithJSON(w, http.StatusOK, response.Response{
		Timestamp: time.Now(),
	})
}

func (h *eventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	claims := r.Context().Value("user").(jwt.MapClaims)
	userID := claims["user_id"].(string)

	err := h.eventService.DeleteEvent(eventID, userID)
	if err != nil {
		switch err.Error(){
		case "[FAIL] event not found":
			respondWithJSON(w, http.StatusNotFound, response.Response{
				Timestamp: time.Now(),
				Message: fmt.Sprintf("%s. Not found", err.Error()),
			})
		case "[FAIL] unauthorized to delete this event":
			respondWithJSON(w, http.StatusForbidden, response.Response{
				Timestamp: time.Now(),
				Message: fmt.Sprintf("%s. Forbidden", err.Error()),
			})
		default:
			respondWithJSON(w, http.StatusInternalServerError, response.Response{
				Timestamp: time.Now(),
				Message: fmt.Sprintf("%s. Internal server error", err.Error()),
			})
		}
	}

	respondWithJSON(w, http.StatusNoContent, response.Response{
		Timestamp: time.Now(),
	})
}

func (h *eventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	event, err := h.eventService.GetEvent(eventID)
	if err != nil {
		if err.Error() == "[FAIL] event not found" {
			respondWithJSON(w, http.StatusNotFound, err.Error())
			return
		}
		respondWithJSON(w, http.StatusInternalServerError, response.Response{
			Timestamp: time.Now(),
			Message: fmt.Sprintf("[FAIL] %s. Internal server error", err.Error()),
		})
	}

	respondWithJSON(w, http.StatusOK, response.Response{
		Timestamp: time.Now(),
		Data: event,
	})
}

func (h *eventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	input := &service.ListEventsInput{
		Page: page,
		PageSize: pageSize,
	}

	events, err := h.eventService.ListEvents(input)
	if err != nil {
		respondWithJSON(w, http.StatusInternalServerError, response.Response{
			Timestamp: time.Now(),
			Message: fmt.Sprintf("[FAIL] %s. Internal server error", err.Error()),
		})
	}

	respondWithJSON(w, http.StatusOK, response.Response{
		Timestamp: time.Now(),
		Data: events,
	})
}
