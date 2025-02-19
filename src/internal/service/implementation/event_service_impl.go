package serviceImplementation

import (
	"context"
	"errors"
	"time"

	"github.com/hafiztri123/src/internal/model"
	"github.com/hafiztri123/src/internal/repository"
	"github.com/hafiztri123/src/internal/service"
)

type eventService struct {
	eventRepository repository.EventRepository 
}

func NewEventService(eventRepo repository.EventRepository) service.EventService {
	return &eventService{
		eventRepository: eventRepo,
	}
}

func (s *eventService) CreateEvent(input *service.CreateEventInput, creatorID string) error {
	event := &model.Event{
		Title: input.Title,
		Description: input.Description,
		StartDate: input.StartDate,
		EndDate: input.EndDate,
		CreatorID: creatorID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := s.eventRepository.Create(context.Background(),event); err != nil {
		return err
	}

	return nil
}

func (s *eventService) UpdateEvent(id string, input *service.UpdateEventInput, userID string)  error {
	event, err := s.eventRepository.GetByID(context.Background(), id)
	if err != nil {
		return err
	}

	if event == nil {
		return errors.New("[FAIL] event not found")
	}

	if event.CreatorID != userID {
		return errors.New("[FAIL] unauthorized to modify this event")
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

func (s *eventService) DeleteEvent(id string, userID string) error {
	event, err := s.eventRepository.GetByID(context.Background(), id)

	if err != nil {
		return err
	}

	if event == nil {
		return errors.New("[FAIL] event not found")
	}

	if event.CreatorID != userID {
		return errors.New("[FAIL] unauthorized to delete this event")
	}

	return s.eventRepository.Delete(context.Background(), id)
}

func (s *eventService) GetEvent(id string) (*model.Event, error) {
	event, err := s.eventRepository.GetByID(context.Background(), id)
	if err != nil {
		return nil, err
	}

	if event == nil {
		return nil, errors.New("[FAIL] event not found")
	}

	return event, nil
}

func (s *eventService) ListEvents(input *service.ListEventsInput) ([]*model.Event, error) {
	offset := (input.Page - 1) * input.PageSize
	return s.eventRepository.List(context.Background(), input.PageSize, offset, input.SortBy, input.SortDir)
}


func (s *eventService) SearchEvents(input *service.SearchEventsInput) (*service.SearchEventsOutput, error){
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

	return &service.SearchEventsOutput{
		Events: events,
		TotalCount: totalCount,
		Page: input.Page,
		PageSize: input.PageSize,
		TotalPages: totalPages,
	}, nil
}
