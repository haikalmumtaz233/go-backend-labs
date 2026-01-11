package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// --- CATEGORY HANDLERS ---

func createCategory(c *gin.Context) {
	var category Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if result := db.Create(&category); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, category)
}

func getCategories(c *gin.Context) {
	var categories []Category
	db.Find(&categories)
	c.JSON(http.StatusOK, categories)
}

// --- SUPPLIER HANDLERS ---

func createSupplier(c *gin.Context) {
	var supplier Supplier
	if err := c.ShouldBindJSON(&supplier); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	if result := db.Create(&supplier); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, supplier)
}

func getSuppliers(c *gin.Context) {
	var suppliers []Supplier
	db.Find(&suppliers)
	c.JSON(http.StatusOK, suppliers)
}

// --- PRODUCT HANDLERS ---

func listProducts(c *gin.Context) {
	var products []Product

	if result := db.Preload("Category").Preload("Supplier").Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func getProductbyId(c *gin.Context) {
	id := c.Param("id")
	var product Product

	if err := db.Preload("Category").Preload("Supplier").First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

func searchProduct(c *gin.Context) {
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Please provide params for search (q)"})
		return
	}

	var products []Product
	searchTerm := "%" + query + "%"

	if result := db.Where("name ILIKE ?", searchTerm).Find(&products); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func createProduct(c *gin.Context) {
	var newProduct Product

	if err := c.ShouldBindJSON(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if newProduct.Name == "" || newProduct.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name cannot be empty, Price cannot be negative"})
		return
	}

	if result := db.Create(&newProduct); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "Success",
		"data":   newProduct,
	})
}

func updateProduct(c *gin.Context) {
	id := c.Param("id")

	var product Product
	var inputProduct Product

	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	if err := c.ShouldBindJSON(&inputProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	db.Model(&product).Updates(inputProduct)

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Data successfully updated",
		"data":    product,
	})
}

func deleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product Product

	if err := db.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Data not found"})
		return
	}

	db.Delete(&product)

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Data successfully deleted",
	})
}

// --- STOCK & HISTORY HANDLERS ---

func adjustStock(c *gin.Context) {
	id := c.Param("id")

	var stockChange StockChange
	if err := c.ShouldBindJSON(&stockChange); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON Body"})
		return
	}

	err := db.Transaction(func(tx *gorm.DB) error {
		var product Product
		if err := tx.First(&product, id).Error; err != nil {
			return err
		}

		switch stockChange.Type {
		case "in":
			product.Stock += stockChange.Amount
		case "out":
			if product.Stock < stockChange.Amount {
				return fmt.Errorf("insufficient stock")
			}
			product.Stock -= stockChange.Amount
		default:
			return fmt.Errorf("invalid type (use 'in' or 'out')")
		}

		if err := tx.Save(&product).Error; err != nil {
			return err
		}

		mutation := StockMutation{
			ProductID: product.ID,
			Amount:    stockChange.Amount,
			Type:      stockChange.Type,
		}

		if err := tx.Create(&mutation).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Transaction Failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Stock adjusted & Logged successfully",
	})
}

func getProductHistory(c *gin.Context) {
	id := c.Param("id")
	var mutations []StockMutation

	if result := db.Where("product_id = ?", id).Order("created_at desc").Find(&mutations); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, mutations)
}
