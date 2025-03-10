package service

import (
	"context"
	"mime/multipart"
	"time"

	"github.com/hafiztri123/src/internal/model"
	errs "github.com/hafiztri123/src/internal/pkg/error"
	"github.com/hafiztri123/src/internal/pkg/storage"
	"github.com/hafiztri123/src/internal/repository"
)

// EventService defines the interface for event-related service operations.
type EventService interface {
    CreateEvent(input *model.CreateEventInput, creatorID string) error
    UpdateEvent(id string, input *model.UpdateEventInput, userID string) error
    DeleteEvent(id string, userID string) error
    GetEvent(id string) (*model.Event, error)
    ListEvents(input *model.ListEventsInput) ([]*model.Event, error)
    SearchEvents(input *model.SearchEventsInput) (*model.SearchEventsOutput, error)
    UploadFile(ctx context.Context,  file multipart.File,input model.UploadFile , eventID string) error
}

// eventService implements the EventService interface.
type eventService struct {
    eventRepository repository.EventRepository
    categoryRepository repository.CategoryRepository
    cloudinary storage.StorageService
}

// NewEventService creates a new instance of EventService.
func NewEventService(eventRepo repository.EventRepository, categoryRepo repository.CategoryRepository, cloudinary storage.StorageService) EventService {
    return &eventService{
        eventRepository: eventRepo,
        categoryRepository: categoryRepo,
        cloudinary: cloudinary,

    }
}

// CreateEvent creates a new event using the provided input and creator ID.
func (s *eventService) CreateEvent(input *model.CreateEventInput, creatorID string) error {
    bool := s.categoryRepository.IsIDExists(input.CategoryID)


    if !bool {
        return errs.NewNotFoundError("Category not found")
    }

    event := &model.Event{
        Title:          input.Title,
        Description:    input.Description,
        CategoryID:     input.CategoryID,
        StartDate:      input.StartDate,
        EndDate:        input.EndDate,
        CreatorID:      creatorID,
        
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    if err := s.eventRepository.Create(context.Background(), event); err != nil {
        return err
    }
    return nil
}

// UpdateEvent updates an existing event if the user is authorized.
func (s *eventService) UpdateEvent(id string, input *model.UpdateEventInput, userID string) error {
    event, err := s.eventRepository.GetByID(context.Background(), id)
    if err != nil {
        return err
    }
    if event == nil {
        return errs.NewNotFoundError("Event not found")
    }
    if event.CreatorID != userID {
        return errs.NewForbiddenError("Creator ID does not matched the required value")
    }
    event.Title = input.Title
    event.Description = input.Description
    event.StartDate = input.StartDate
    event.EndDate = input.EndDate
    event.UpdatedAt = time.Now()
    if err := s.eventRepository.Update(context.Background(), event); err != nil {
        return err
    }
    return nil
}

// DeleteEvent deletes an event if the user is authorized.
func (s *eventService) DeleteEvent(id string, userID string) error {
    event, err := s.eventRepository.GetByID(context.Background(), id)
    if err != nil {
        return err
    }
    if event == nil {
        return errs.NewNotFoundError("Event not found")
    }
    if event.CreatorID != userID {
        return errs.NewForbiddenError("Creator ID does not matched the required value")
    }
    return s.eventRepository.Delete(context.Background(), id)
}

// GetEvent retrieves an event by its ID.
func (s *eventService) GetEvent(id string) (*model.Event, error) {
    event, err := s.eventRepository.GetByID(context.Background(), id)
    if err != nil {
        return nil, err
    }
    if event == nil {
        return nil, errs.NewNotFoundError("Event not found")
    }
    return event, nil
}

// ListEvents retrieves a paginated list of events based on the input parameters.
func (s *eventService) ListEvents(input *model.ListEventsInput) ([]*model.Event, error) {
    offset := (input.Page - 1) * input.PageSize
    return s.eventRepository.List(context.Background(), input.PageSize, offset, input.SortBy, input.SortDir)
}

// SearchEvents searches for events based on the input parameters and returns paginated results.
func (s *eventService) SearchEvents(input *model.SearchEventsInput) (*model.SearchEventsOutput, error) {
    if input.Page < 1 {
        input.Page = 1
    }

    if input.PageSize < 1 || input.PageSize > 100 {
        input.PageSize = 10
    }

    events, totalCount, err := s.eventRepository.Search(context.Background(), input)
    if err != nil {
        return nil, err
    }

    totalPages := int(totalCount) / input.PageSize
    if int(totalCount)%input.PageSize > 0 {
        totalPages++
    }

    return &model.SearchEventsOutput{
        Events:     events,
        TotalCount: totalCount,
        Page:       input.Page,
        PageSize:   input.PageSize,
        TotalPages: totalPages,
    }, nil
}


func (s *eventService ) UploadFile(ctx context.Context, file multipart.File,input model.UploadFile ,eventID string) error {

    fileURL, err := s.cloudinary.UploadFile(ctx, file, input.FileName)
    if err != nil {
        return err
    }

    fileModel := &model.File{
        EventID: eventID,
        FileName: input.FileName,
        FileType: input.FileType,
        FileURL: fileURL,
        CreatedAt: time.Now(),
    }

    err = s.eventRepository.UploadFile(ctx, fileModel)
    if err != nil {
        return err
    }

    return nil


}
