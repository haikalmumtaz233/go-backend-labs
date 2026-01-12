package entity

import (
	"time"
)

type Event struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Title            string    `gorm:"type:varchar(200);not null" json:"title"`
	Description      string    `gorm:"type:text" json:"description"`
	Date             time.Time `gorm:"not null" json:"date"`
	Location         string    `gorm:"type:varchar(200);not null" json:"location"`
	TotalTickets     int       `gorm:"not null" json:"total_tickets"`
	AvailableTickets int       `gorm:"not null" json:"available_tickets"`
	Price            float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	CreatedAt        time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type CreateEventInput struct {
	Title        string    `json:"title" binding:"required,min=3,max=200"`
	Description  string    `json:"description"`
	Date         time.Time `json:"date" binding:"required"`
	Location     string    `json:"location" binding:"required"`
	TotalTickets int       `json:"total_tickets" binding:"required,min=1"`
	Price        float64   `json:"price" binding:"required,min=0"`
}

type UpdateEventInput struct {
	Title        string    `json:"title" binding:"omitempty,min=3,max=200"`
	Description  string    `json:"description"`
	Date         time.Time `json:"date"`
	Location     string    `json:"location"`
	TotalTickets int       `json:"total_tickets" binding:"omitempty,min=1"`
	Price        float64   `json:"price" binding:"omitempty,min=0"`
}

type EventFilter struct {
	Search   string
	Location string
	DateFrom time.Time
	DateTo   time.Time
	Page     int
	PageSize int
}
