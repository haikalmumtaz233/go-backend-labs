package main

import (
	"fmt"
	"net/http"
)

func main() {
	ConnectDB()
	mux := http.NewServeMux()

	// Product Routes
	mux.HandleFunc("GET /products", listProducts)
	mux.HandleFunc("GET /products/{product_id}", getProductbyId)
	mux.HandleFunc("GET /products/search", searchProduct)
	mux.HandleFunc("POST /products", createProduct)
	mux.HandleFunc("PUT /products/{product_id}", updateProduct)
	mux.HandleFunc("PATCH /products/{product_id}/stock", adjustStock)
	mux.HandleFunc("GET /products/{product_id}/history", getProductHistory)
	mux.HandleFunc("DELETE /products/{product_id}", deleteProduct)

	// Category Routes
	mux.HandleFunc("POST /categories", createCategory)
	mux.HandleFunc("GET /categories", getCategories)

	// Supplier Routes
	mux.HandleFunc("POST /suppliers", createSupplier)
	mux.HandleFunc("GET /suppliers", getSuppliers)

	fmt.Println("Server running on http://localhost:8080...")
	http.ListenAndServe(":8080", mux)
}
