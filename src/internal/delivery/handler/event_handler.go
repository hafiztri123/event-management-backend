package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"github.com/hafiztri123/src/internal/pkg/response"
	"github.com/hafiztri123/src/internal/service"
)

// EventHandler defines the interface for event-related HTTP handlers.
type EventHandler interface {
    CreateEvent(w http.ResponseWriter, r *http.Request)
    UpdateEvent(w http.ResponseWriter, r *http.Request)
    DeleteEvent(w http.ResponseWriter, r *http.Request)
    GetEvent(w http.ResponseWriter, r *http.Request)
    ListEvents(w http.ResponseWriter, r *http.Request)
    SearchEvents(w http.ResponseWriter, r *http.Request)
    UploadFile(w http.ResponseWriter, r *http.Request)
}

// eventHandler implements the EventHandler interface.
type eventHandler struct {
    eventService service.EventService
    validator    *validator.Validate
}

// NewEventHandler creates a new instance of EventHandler.
func NewEventHandler(eventService service.EventService) EventHandler {
    return &eventHandler{
        eventService: eventService,
        validator:    validator.New(),
    }
}

// CreateEvent godoc
// @Summary      Create new event
// @Description  Create a new event with the given details
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        input body service.CreateEventInput true "Event Details"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     Bearer
// @Router       /events [post]
func (h *eventHandler) CreateEvent(w http.ResponseWriter, r *http.Request) {
    var input model.CreateEventInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        HandleErrorResponse(w, err)
        return
    }

    if err := h.validator.Struct(input); err != nil {
        HandleErrorResponse(w, err)
        return
    }

    claims := r.Context().Value("user").(jwt.MapClaims)
    userID := claims["user_id"].(string)

    err := h.eventService.CreateEvent(&input, userID)
    if err != nil {
        HandleErrorResponse(w, err)
        return
    }

    respondWithJSON(w, http.StatusCreated, response.Response{
        Timestamp: time.Now(),
    })
}

// UpdateEvent godoc
// @Summary      Update event
// @Description  Update an existing event with new details
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Event ID"
// @Param        input body service.UpdateEventInput true "Event Details"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     Bearer
// @Router       /events/{id} [put]
func (h *eventHandler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
    eventID := chi.URLParam(r, "id")
    var input model.UpdateEventInput
    if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
        HandleErrorResponse(w, err)
        return
    }

    if err := h.validator.Struct(input); err != nil {
        HandleErrorResponse(w, errs.NewBadRequestError("Request is not valid"))
        return
    }

    claims := r.Context().Value("user").(jwt.MapClaims)
    userID := claims["user_id"].(string)

    err := h.eventService.UpdateEvent(eventID, &input, userID)
    if err != nil {
        HandleErrorResponse(w, err)
        return
    }

    respondWithJSON(w, http.StatusOK, response.Response{
        Timestamp: time.Now(),
    })
}

// DeleteEvent godoc
// @Summary      Delete event
// @Description  Delete an existing event
// @Tags         events
// @Produce      json
// @Param        id   path      string  true  "Event ID"
// @Success      204  {object}  response.Response
// @Failure      401  {object}  response.Response
// @Failure      403  {object}  response.Response
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Security     Bearer
// @Router       /events/{id} [delete]
func (h *eventHandler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
    eventID := chi.URLParam(r, "id")

    claims := r.Context().Value("user").(jwt.MapClaims)
    userID := claims["user_id"].(string)

    err := h.eventService.DeleteEvent(eventID, userID)
    if err != nil {
        HandleErrorResponse(w, err)
        return
    }

    respondWithJSON(w, http.StatusNoContent, response.Response{
        Timestamp: time.Now(),
    })
}

// GetEvent godoc
// @Summary      Get event details
// @Description  Get details of a specific event
// @Tags         events
// @Produce      json
// @Param        id   path      string  true  "Event ID"
// @Success      200  {object}  response.Response{data=model.Event}
// @Failure      404  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /events/{id} [get]
func (h *eventHandler) GetEvent(w http.ResponseWriter, r *http.Request) {
    eventID := chi.URLParam(r, "id")

    event, err := h.eventService.GetEvent(eventID)
    if err != nil {
        HandleErrorResponse(w, err)
        return
    }

    respondWithJSON(w, http.StatusOK, response.Response{
        Timestamp: time.Now(),
        Data:      event,
    })
}

// ListEvents godoc
// @Summary      List events
// @Description  Get a paginated list of all events
// @Tags         events
// @Produce      json
// @Param        page      query     int  false  "Page number"  minimum(1)
// @Param        page_size query     int  false  "Page size"    minimum(1)  maximum(100)
// @Success      200  {object}  response.Response{data=[]model.Event}
// @Failure      500  {object}  response.Response
// @Router       /events [get]
func (h *eventHandler) ListEvents(w http.ResponseWriter, r *http.Request) {
    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }

    pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
    if err != nil || pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }

    input := &model.ListEventsInput{
        Page:     page,
        PageSize: pageSize,
    }

    events, err := h.eventService.ListEvents(input)
    if err != nil {
        HandleErrorResponse(w, err)
        return
    }

    respondWithJSON(w, http.StatusOK, response.Response{
        Timestamp: time.Now(),
        Data:      events,
    })
}

// SearchEvents godoc
// @Summary      Search events
// @Description  Search and filter events with various criteria
// @Tags         events
// @Accept       json
// @Produce      json
// @Param        query      query     string  false  "Search term in title and description"
// @Param        start_date query     string  false  "Filter events starting after this date (RFC3339)"
// @Param        end_date   query     string  false  "Filter events ending before this date (RFC3339)"
// @Param        creator    query     string  false  "Filter by creator ID"
// @Param        page       query     int     false  "Page number"  minimum(1)
// @Param        page_size  query     int     false  "Page size"    minimum(1)  maximum(100)
// @Param        sort_by    query     string  false  "Sort field (title, start_date, end_date, created_at)"
// @Param        sort_dir   query     string  false  "Sort direction (asc, desc)"
// @Success      200  {object}  response.Response{data=service.SearchEventsOutput}
// @Failure      400  {object}  response.Response
// @Failure      500  {object}  response.Response
// @Router       /events/search [get]
func (h *eventHandler) SearchEvents(w http.ResponseWriter, r *http.Request) {
    query := r.URL.Query().Get("query")
    creator := r.URL.Query().Get("creator")

    var startDate, endDate *time.Time

    if startDateStr := r.URL.Query().Get("start_date"); startDateStr != "" {
        parsedDate, err := time.Parse(time.RFC3339, startDateStr)
        if err == nil {
            startDate = &parsedDate
        }
    }

    if endDateStr := r.URL.Query().Get("end_date"); endDateStr != "" {
        parsedDate, err := time.Parse(time.RFC3339, endDateStr)
        if err == nil {
            endDate = &parsedDate
        }
    }

    page, err := strconv.Atoi(r.URL.Query().Get("page"))
    if err != nil || page < 1 {
        page = 1
    }

    pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
    if err != nil || pageSize < 1 || pageSize > 100 {
        pageSize = 10
    }

    sortBy := r.URL.Query().Get("sort_by")
    sortDir := r.URL.Query().Get("sort_dir")

    input := &model.SearchEventsInput{
        Query:    query,
        StartDate: startDate,
        EndDate:   endDate,
        Creator:   creator,
        Page:      page,
        PageSize:  pageSize,
        SortBy:    sortBy,
        SortDir:   sortDir,
    }

    result, err := h.eventService.SearchEvents(input)
    if err != nil {
        HandleErrorResponse(w, err)
        return
    }

    respondWithJSON(w, http.StatusOK, response.Response{
        Timestamp: time.Now(),
        Data:      result,
    })
}

func (h *eventHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
    eventID := chi.URLParam(r, "id")

    if err := r.ParseMultipartForm(10 << 20); err != nil {
		HandleErrorResponse(w, errs.NewBadRequestError("File too large"))
		return
	}

    file, header, err := r.FormFile("file")
    if err != nil {
		HandleErrorResponse(w, errs.NewBadRequestError("Invalid file"))
		return
	}
	defer file.Close()

	if header.Size > int64(10 << 20) { 
        HandleErrorResponse(w, errs.NewEntityTooLargeError("File too large"))
        return
    }

    if err := h.eventService.UploadFile(r.Context(), file, model.UploadFile{
        FileName: header.Filename,
        FileType: header.Header.Get("Content-Type"),
    }, eventID); err != nil {
        HandleErrorResponse(w, err)
    }

    respondWithJSON(w, http.StatusOK, response.Response{
        Timestamp: time.Now(),
    })


}


