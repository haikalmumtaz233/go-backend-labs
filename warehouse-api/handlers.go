package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gorm.io/gorm"
)

// --- CATEGORY HANDLERS ---

func createCategory(w http.ResponseWriter, r *http.Request) {
	var category Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	db.Create(&category)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func getCategories(w http.ResponseWriter, r *http.Request) {
	var categories []Category
	db.Find(&categories)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}

// --- SUPPLIER HANDLERS ---

func createSupplier(w http.ResponseWriter, r *http.Request) {
	var supplier Supplier
	if err := json.NewDecoder(r.Body).Decode(&supplier); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	db.Create(&supplier)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(supplier)
}

func getSuppliers(w http.ResponseWriter, r *http.Request) {
	var suppliers []Supplier
	db.Find(&suppliers)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(suppliers)
}

// --- PRODUCT HANDLERS ---

func listProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product

	if result := db.Preload("Category").Preload("Supplier").Find(&products); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProductbyId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("product_id")
	var product Product

	if err := db.Preload("Category").Preload("Supplier").First(&product, id).Error; err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func searchProduct(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")

	if query == "" {
		http.Error(w, "Please provide params for search", http.StatusBadRequest)
		return
	}

	var products []Product
	searchTerm := "%" + query + "%"

	if result := db.Where("name ILIKE ?", searchTerm).Find(&products); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func createProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct Product

	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	if newProduct.Name == "" || newProduct.Price < 0 {
		http.Error(w, "Nama tidak boleh kosong, price tidak boleh dibawah 0!", http.StatusBadRequest)
		return
	}

	if result := db.Create(&newProduct); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "Success",
		"data":   newProduct,
	})
}

func updateProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("product_id")

	var product Product
	var inputProduct Product

	if err := db.First(&product, id).Error; err != nil {
		http.Error(w, "Data not Found", http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&inputProduct); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		return
	}

	db.Model(&product).Updates(inputProduct)

	response := map[string]interface{}{
		"status":  "Success",
		"message": "Data berhasil diupdate",
		"data":    product,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func deleteProduct(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("product_id")
	var product Product

	if err := db.First(&product, id).Error; err != nil {
		http.Error(w, "Data not Found", http.StatusNotFound)
		return
	}

	db.Delete(&product)

	response := map[string]interface{}{
		"status":  "Success",
		"message": "Data berhasil dihapus",
		"data":    product,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func adjustStock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("product_id")

	var stockChange StockChange
	if err := json.NewDecoder(r.Body).Decode(&stockChange); err != nil {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
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
			return fmt.Errorf("invalid type")
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
		http.Error(w, "Transaction Failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"status":  "Success",
		"message": "Stock adjusted & Logged successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getProductHistory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("product_id")
	var mutations []StockMutation

	if result := db.Where("product_id = ?", id).Order("created_at desc").Find(&mutations); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mutations)
}
