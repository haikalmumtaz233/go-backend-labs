package main

import (
	"github.com/haikalmumtaz233/warehouse-gin/database"
	"github.com/haikalmumtaz233/warehouse-gin/handler"
	"github.com/haikalmumtaz233/warehouse-gin/repository"
	"github.com/haikalmumtaz233/warehouse-gin/service"

	"github.com/gin-gonic/gin"
)

func main() {
	db := database.ConnectDB()

	productRepo := repository.NewProductRepository(db)
	productService := service.NewProductService(productRepo)
	productHandler := handler.NewProductHandler(productService)

	r := gin.Default()
	api := r.Group("/api")
	{
		api.GET("/products", productHandler.GetProducts)
		api.GET("/products/:id", productHandler.GetProductByID)
		api.POST("/products", productHandler.CreateProduct)
		api.PUT("/products/:id", productHandler.UpdateProduct)
		api.DELETE("/products/:id", productHandler.DeleteProduct)

		api.PATCH("/products/:id/stock", productHandler.AdjustStock)
		api.GET("/products/:id/history", productHandler.GetHistory)
	}

	r.Run(":8080")
}
