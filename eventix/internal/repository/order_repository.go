package repository

import (
	"eventix/internal/entity"

	"gorm.io/gorm"
)

type OrderRepository interface {
	Save(tx *gorm.DB, order *entity.Order) error
	FindByID(id uint) (*entity.Order, error)
	FindByUserID(userID uint) ([]entity.Order, error)
	UpdateStatus(tx *gorm.DB, orderID uint, status entity.OrderStatus) error
	GetDB() *gorm.DB
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Save(tx *gorm.DB, order *entity.Order) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Create(order).Error
}

func (r *orderRepository) FindByID(id uint) (*entity.Order, error) {
	var order entity.Order
	if err := r.db.Preload("Event").Preload("Tickets").First(&order, id).Error; err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindByUserID(userID uint) ([]entity.Order, error) {
	var orders []entity.Order
	if err := r.db.Preload("Event").Preload("Tickets").Where("user_id = ?", userID).Order("created_at DESC").Find(&orders).Error; err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *orderRepository) UpdateStatus(tx *gorm.DB, orderID uint, status entity.OrderStatus) error {
	if tx == nil {
		tx = r.db
	}
	return tx.Model(&entity.Order{}).Where("id = ?", orderID).Update("status", status).Error
}

func (r *orderRepository) GetDB() *gorm.DB {
	return r.db
}
