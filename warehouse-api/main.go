package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model

	Name  string `json:"product_name" gorm:"type:varchar(100);not null"`
	Price int    `json:"product_price" gorm:"not null"`
	Stock int    `json:"product_stock" gorm:"not null"`
}

type StockChange struct {
	Amount int    `json:"amount"`
	Type   string `json:"type"`
}

var db *gorm.DB

func ConnectDB() {
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		dbHost, dbUser, dbPass, dbName, dbPort)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the Warehouse Database!", err.Error())
	}

	db.AutoMigrate(&Product{})
	fmt.Println("Successfully Connected to the Warehouse Database!")
}

// HANDLERS

func listProducts(w http.ResponseWriter, r *http.Request) {
	var products []Product

	if result := db.Find(&products); result.Error != nil {
		http.Error(w, result.Error.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func getProductbyId(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("product_id")
	var product Product

	if err := db.First(&product, id).Error; err != nil {
		http.Error(w, "Data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
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

func adjustStock(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("product_id")

	var product Product
	if err := db.First(&product, id).Error; err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	var stockChange StockChange
	if err := json.NewDecoder(r.Body).Decode(&stockChange); err != nil {
		http.Error(w, "Invalid JSON Body", http.StatusBadRequest)
		return
	}

	switch stockChange.Type {
	case "in":
		product.Stock += stockChange.Amount
	case "out":
		if product.Stock < stockChange.Amount {
			http.Error(w, "Insufficient stock!", http.StatusBadRequest)
			return
		}
		product.Stock -= stockChange.Amount
	default:
		http.Error(w, "Invalid stock type. Use 'in' or 'out'", http.StatusBadRequest)
		return
	}

	if err := db.Save(&product).Error; err != nil {
		http.Error(w, "Failed to update stock", http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status":     "Success",
		"message":    "Stock adjusted succesfully",
		"new_stock":  product.Stock,
		"adjustment": stockChange,
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

func main() {

	ConnectDB()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /products", listProducts)
	mux.HandleFunc("GET /products/{product_id}", getProductbyId)
	mux.HandleFunc("GET /products/search", searchProduct)

	mux.HandleFunc("POST /products", createProduct)

	mux.HandleFunc("PUT /products/{product_id}", updateProduct)
	mux.HandleFunc("PATCH /products/{product_id}/stock", adjustStock)

	mux.HandleFunc("DELETE /products/{product_id}", deleteProduct)

	fmt.Println("Server running on http://localhost:8080...")
	http.ListenAndServe(":8080", mux)
}
