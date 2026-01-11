package repository

import (
	"errors"

	"github.com/haikalmumtaz233/warehouse-gin/entity"

	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll() ([]entity.Product, error)
	FindByID(id string) (entity.Product, error)
	Search(name string) ([]entity.Product, error)
	Save(product entity.Product) (entity.Product, error)
	Update(product entity.Product) (entity.Product, error)
	Delete(product entity.Product) error

	AdjustStock(productID string, amount int, typeChange string) error
	GetHistory(productID string) ([]entity.StockMutation, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) FindAll() ([]entity.Product, error) {
	var products []entity.Product
	err := r.db.Preload("Category").Preload("Supplier").Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id string) (entity.Product, error) {
	var product entity.Product
	err := r.db.Preload("Category").Preload("Supplier").First(&product, id).Error
	return product, err
}

func (r *productRepository) Search(name string) ([]entity.Product, error) {
	var products []entity.Product
	searchTerm := "%" + name + "%"
	err := r.db.Where("name ILIKE ?", searchTerm).Find(&products).Error
	return products, err
}

func (r *productRepository) Save(product entity.Product) (entity.Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

func (r *productRepository) Update(product entity.Product) (entity.Product, error) {
	err := r.db.Save(&product).Error
	return product, err
}

func (r *productRepository) Delete(product entity.Product) error {
	return r.db.Delete(&product).Error
}

func (r *productRepository) GetHistory(productID string) ([]entity.StockMutation, error) {
	var mutations []entity.StockMutation
	err := r.db.Where("product_id = ?", productID).Order("created_at desc").Find(&mutations).Error
	return mutations, err
}

func (r *productRepository) AdjustStock(productID string, amount int, typeChange string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var product entity.Product
		if err := tx.First(&product, productID).Error; err != nil {
			return err
		}

		if typeChange == "in" {
			product.Stock += amount
		} else if typeChange == "out" {
			if product.Stock < amount {
				return errors.New("insufficient stock")
			}
			product.Stock -= amount
		}

		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		mutation := entity.StockMutation{
			ProductID: product.ID,
			Amount:    amount,
			Type:      typeChange,
		}
		if err := tx.Create(&mutation).Error; err != nil {
			return err
		}

		return nil
	})
}
