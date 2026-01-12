package repository

import (
	"eventix/internal/entity"

	"gorm.io/gorm"
)

type TicketRepository interface {
	SaveBatch(tx *gorm.DB, tickets []entity.Ticket) error
	FindByOrderID(orderID uint) ([]entity.Ticket, error)
	FindByTicketCode(code string) (*entity.Ticket, error)
	UpdateStatus(ticketID uint, status entity.TicketStatus) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db: db}
}

func (r *ticketRepository) SaveBatch(tx *gorm.DB, tickets []entity.Ticket) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(&tickets).Error
}

func (r *ticketRepository) FindByOrderID(orderID uint) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	if err := r.db.Where("order_id = ?", orderID).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) FindByTicketCode(code string) (*entity.Ticket, error) {
	var ticket entity.Ticket
	if err := r.db.Preload("Event").Where("ticket_code = ?", code).First(&ticket).Error; err != nil {
		return nil, err
	}
	return &ticket, nil
}

func (r *ticketRepository) UpdateStatus(ticketID uint, status entity.TicketStatus) error {
	return r.db.Model(&entity.Ticket{}).Where("id = ?", ticketID).Update("status", status).Error
}
