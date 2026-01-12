package entity

import (
	"time"
)

type TicketStatus string

const (
	TicketStatusValid TicketStatus = "VALID"
	TicketStatusUsed  TicketStatus = "USED"
)

type Ticket struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	OrderID    uint         `gorm:"not null;index" json:"order_id"`
	EventID    uint         `gorm:"not null;index" json:"event_id"`
	TicketCode string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"ticket_code"`
	Status     TicketStatus `gorm:"type:varchar(20);default:'VALID'" json:"status"`
	CreatedAt  time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time    `gorm:"autoUpdateTime" json:"updated_at"`

	Order Order `gorm:"foreignKey:OrderID" json:"-"`
	Event Event `gorm:"foreignKey:EventID" json:"event,omitempty"`
}
