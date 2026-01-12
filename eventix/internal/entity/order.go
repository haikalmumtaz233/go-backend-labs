package entity

import (
	"time"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "PENDING"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

type Order struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	UserID      uint        `gorm:"not null;index" json:"user_id"`
	EventID     uint        `gorm:"not null;index" json:"event_id"`
	Quantity    int         `gorm:"not null" json:"quantity"`
	TotalAmount float64     `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      OrderStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	CreatedAt   time.Time   `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time   `gorm:"autoUpdateTime" json:"updated_at"`

	Event   Event    `gorm:"foreignKey:EventID" json:"event,omitempty"`
	User    User     `gorm:"foreignKey:UserID" json:"-"`
	Tickets []Ticket `gorm:"foreignKey:OrderID" json:"tickets,omitempty"`
}

type BookingInput struct {
	Quantity int `json:"qty" binding:"required,min=1,max=10"`
}

type PaymentInput struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}
