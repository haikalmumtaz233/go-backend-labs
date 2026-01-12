package service

import (
	"errors"

	"eventix/internal/entity"
	"eventix/internal/repository"
)

var (
	ErrEventNotFound = errors.New("event not found")
)

type EventService interface {
	CreateEvent(input *entity.CreateEventInput) (*entity.Event, error)
	GetAllEvents(filter entity.EventFilter) ([]entity.Event, int64, error)
	GetEventByID(id uint) (*entity.Event, error)
	UpdateEvent(id uint, input *entity.UpdateEventInput) (*entity.Event, error)
	DeleteEvent(id uint) error
}

type eventService struct {
	eventRepo repository.EventRepository
}

func NewEventService(eventRepo repository.EventRepository) EventService {
	return &eventService{eventRepo: eventRepo}
}

func (s *eventService) CreateEvent(input *entity.CreateEventInput) (*entity.Event, error) {
	event := &entity.Event{
		Title:            input.Title,
		Description:      input.Description,
		Date:             input.Date,
		Location:         input.Location,
		TotalTickets:     input.TotalTickets,
		AvailableTickets: input.TotalTickets,
		Price:            input.Price,
	}

	if err := s.eventRepo.Save(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventService) GetAllEvents(filter entity.EventFilter) ([]entity.Event, int64, error) {
	return s.eventRepo.FindAll(filter)
}

func (s *eventService) GetEventByID(id uint) (*entity.Event, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, ErrEventNotFound
	}
	return event, nil
}

func (s *eventService) UpdateEvent(id uint, input *entity.UpdateEventInput) (*entity.Event, error) {
	event, err := s.eventRepo.FindByID(id)
	if err != nil {
		return nil, ErrEventNotFound
	}

	if input.Title != "" {
		event.Title = input.Title
	}
	if input.Description != "" {
		event.Description = input.Description
	}
	if !input.Date.IsZero() {
		event.Date = input.Date
	}
	if input.Location != "" {
		event.Location = input.Location
	}
	if input.TotalTickets > 0 {
		diff := input.TotalTickets - event.TotalTickets
		event.TotalTickets = input.TotalTickets
		event.AvailableTickets += diff
		if event.AvailableTickets < 0 {
			event.AvailableTickets = 0
		}
	}
	if input.Price > 0 {
		event.Price = input.Price
	}

	if err := s.eventRepo.Update(event); err != nil {
		return nil, err
	}

	return event, nil
}

func (s *eventService) DeleteEvent(id uint) error {
	_, err := s.eventRepo.FindByID(id)
	if err != nil {
		return ErrEventNotFound
	}
	return s.eventRepo.Delete(id)
}
