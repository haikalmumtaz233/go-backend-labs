package service

import (
	"errors"

	"github.com/haikalmumtaz233/warehouse-gin/entity"
	"github.com/haikalmumtaz233/warehouse-gin/repository"
)

type ProductService interface {
	GetAll() ([]entity.Product, error)
	GetByID(id string) (entity.Product, error)
	Search(name string) ([]entity.Product, error)
	Create(input entity.Product) (entity.Product, error)
	Update(id string, input entity.Product) (entity.Product, error)
	Delete(id string) error

	AdjustStock(id string, input entity.StockChangeRequest) error
	GetHistory(id string) ([]entity.StockMutation, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) GetAll() ([]entity.Product, error) {
	return s.repo.FindAll()
}

func (s *productService) GetByID(id string) (entity.Product, error) {
	return s.repo.FindByID(id)
}

func (s *productService) Search(name string) ([]entity.Product, error) {
	return s.repo.Search(name)
}

func (s *productService) Create(input entity.Product) (entity.Product, error) {
	if input.Price < 0 {
		return entity.Product{}, errors.New("price cannot be negative")
	}
	if input.Name == "" {
		return entity.Product{}, errors.New("name is required")
	}
	return s.repo.Save(input)
}

func (s *productService) Update(id string, input entity.Product) (entity.Product, error) {
	existingProduct, err := s.repo.FindByID(id)
	if err != nil {
		return entity.Product{}, err
	}

	existingProduct.Name = input.Name
	existingProduct.Price = input.Price

	return s.repo.Update(existingProduct)
}

func (s *productService) Delete(id string) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(product)
}

func (s *productService) AdjustStock(id string, input entity.StockChangeRequest) error {
	return s.repo.AdjustStock(id, input.Amount, input.Type)
}

func (s *productService) GetHistory(id string) ([]entity.StockMutation, error) {
	return s.repo.GetHistory(id)
}
