package repository

import (
	"eventix/internal/entity"

	"gorm.io/gorm"
)

type EventRepository interface {
	FindAll(filter entity.EventFilter) ([]entity.Event, int64, error)
	FindByID(id uint) (*entity.Event, error)
	Save(event *entity.Event) error
	Update(event *entity.Event) error
	Delete(id uint) error
	DecrementAvailableTickets(tx *gorm.DB, eventID uint, qty int) error
	IncrementAvailableTickets(tx *gorm.DB, eventID uint, qty int) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) EventRepository {
	return &eventRepository{db: db}
}

func (r *eventRepository) FindAll(filter entity.EventFilter) ([]entity.Event, int64, error) {
	var events []entity.Event
	var total int64

	query := r.db.Model(&entity.Event{})

	if filter.Search != "" {
		searchPattern := "%" + filter.Search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", searchPattern, searchPattern)
	}

	if filter.Location != "" {
		query = query.Where("location ILIKE ?", "%"+filter.Location+"%")
	}

	if !filter.DateFrom.IsZero() {
		query = query.Where("date >= ?", filter.DateFrom)
	}

	if !filter.DateTo.IsZero() {
		query = query.Where("date <= ?", filter.DateTo)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 10
	}

	offset := (filter.Page - 1) * filter.PageSize
	if err := query.Offset(offset).Limit(filter.PageSize).Order("date ASC").Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, total, nil
}

func (r *eventRepository) FindByID(id uint) (*entity.Event, error) {
	var event entity.Event
	if err := r.db.First(&event, id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (r *eventRepository) Save(event *entity.Event) error {
	return r.db.Create(event).Error
}

func (r *eventRepository) Update(event *entity.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Event{}, id).Error
}

func (r *eventRepository) DecrementAvailableTickets(tx *gorm.DB, eventID uint, qty int) error {
	return tx.Model(&entity.Event{}).
		Where("id = ? AND available_tickets >= ?", eventID, qty).
		UpdateColumn("available_tickets", gorm.Expr("available_tickets - ?", qty)).Error
}

func (r *eventRepository) IncrementAvailableTickets(tx *gorm.DB, eventID uint, qty int) error {
	return tx.Model(&entity.Event{}).
		Where("id = ?", eventID).
		UpdateColumn("available_tickets", gorm.Expr("available_tickets + ?", qty)).Error
}
