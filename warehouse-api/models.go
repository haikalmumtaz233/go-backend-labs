package main

import "gorm.io/gorm"

type Category struct {
	gorm.Model
	Name string `json:"name" gorm:"type:varchar(100);not null"`
}

type Supplier struct {
	gorm.Model
	Name  string `json:"name" gorm:"type:varchar(100);not null"`
	Email string `json:"email" gorm:"type:varchar(100)"`
	Phone string `json:"phone" gorm:"type:varchar(20)"`
}

type Product struct {
	gorm.Model
	Name  string `json:"product_name" gorm:"type:varchar(100);not null"`
	Price int    `json:"product_price" gorm:"not null"`
	Stock int    `json:"product_stock" gorm:"not null"`

	CategoryID uint     `json:"category_id"`
	Category   Category `json:"category" gorm:"foreignKey:CategoryID"`

	SupplierID uint     `json:"supplier_id"`
	Supplier   Supplier `json:"supplier" gorm:"foreignKey:SupplierID"`
}

type StockChange struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

type StockMutation struct {
	gorm.Model
	ProductID uint    `json:"product_id"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Amount    int     `json:"amount"`
	Type      string  `json:"type"`
}
