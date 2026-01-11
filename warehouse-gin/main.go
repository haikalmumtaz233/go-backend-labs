package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	ConnectDB()

	r := gin.Default()

	api := r.Group("/api")
	{
		// -- Product Routes
		api.GET("/products", listProducts)
		api.GET("/products/:id", getProductbyId)
		api.GET("/products/search", searchProduct)
		api.GET("/products/:id/history", getProductHistory)
		api.POST("/products", createProduct)
		api.PUT("/products/:id", updateProduct)
		api.PATCH("/products/:id/stock", adjustStock)
		api.DELETE("/products/:id", deleteProduct)

		// Category Routes
		api.GET("/categories", getCategories)
		api.POST("/categories", createCategory)

		// Supplier Routes
		api.GET("/suppliers", getSuppliers)
		api.POST("/suppliers", createSupplier)
	}

	fmt.Println("Server running on http://localhost:8080...")
	r.Run(":8080")
}
