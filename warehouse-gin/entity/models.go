package entity

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
}

type Supplier struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `json:"name" gorm:"type:varchar(100);not null"`
	Email     string         `json:"email" gorm:"type:varchar(100)"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
}

type Product struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Name      string         `json:"product_name" gorm:"type:varchar(100);not null"`
	Price     int            `json:"product_price" gorm:"not null"`
	Stock     int            `json:"product_stock" gorm:"not null"`

	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`
	SupplierID uint     `json:"supplier_id"`
	Supplier   Supplier `json:"supplier" gorm:"foreignKey:SupplierID"`
}

type StockMutation struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	ProductID uint      `json:"product_id"`
	Amount    int       `json:"amount"`
	Type      string    `json:"type"`
}

type StockChangeRequest struct {
	Amount int    `json:"amount" binding:"required,min=1"`
	Type   string `json:"type" binding:"required,oneof=in out"`
}
