package main

import (
	"log"
	"os"

	"eventix/internal/entity"
	"eventix/internal/handler"
	"eventix/internal/middleware"
	"eventix/internal/repository"
	"eventix/internal/service"
	"eventix/pkg/database"
	"eventix/pkg/worker"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// ==========================================================
	// Step 1: Load environment variables
	// ==========================================================
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// ==========================================================
	// Step 2: Initialize Database Connection
	// ==========================================================
	db, err := database.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run database migrations for all entities
	if err := database.AutoMigrate(db,
		&entity.User{},
		&entity.Event{},
		&entity.Order{},
		&entity.Ticket{},
	); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// ==========================================================
	// Step 3: Initialize Email Worker Channel and Goroutine
	// ==========================================================
	emailChan := make(chan worker.EmailJob, 100)
	worker.StartEmailWorker(emailChan)

	// ==========================================================
	// Step 4: Dependency Injection - Repositories
	// ==========================================================
	userRepo := repository.NewUserRepository(db)
	eventRepo := repository.NewEventRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	ticketRepo := repository.NewTicketRepository(db)

	// ==========================================================
	// Step 5: Dependency Injection - Services
	// ==========================================================
	authService := service.NewAuthService(userRepo)
	eventService := service.NewEventService(eventRepo)
	orderService := service.NewOrderService(orderRepo, eventRepo, ticketRepo, emailChan)

	// ==========================================================
	// Step 6: Dependency Injection - Handlers
	// ==========================================================
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler()
	eventHandler := handler.NewEventHandler(eventService)
	orderHandler := handler.NewOrderHandler(orderService)

	// ==========================================================
	// Step 7: Setup Gin Router and Routes
	// ==========================================================
	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "eventix",
		})
	})

	// API routes
	api := router.Group("/api")
	{
		// Public authentication routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected user routes
		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/profile", userHandler.GetProfile)
		}

		// Event routes
		events := api.Group("/events")
		{
			// Public event routes
			events.GET("", eventHandler.GetAllEvents)
			events.GET("/:id", eventHandler.GetEventByID)

			// Protected booking route
			events.POST("/:id/book", middleware.AuthMiddleware(), orderHandler.BookTickets)

			// Admin-only event management routes
			adminEvents := events.Group("")
			adminEvents.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
			{
				adminEvents.POST("", eventHandler.CreateEvent)
				adminEvents.PUT("/:id", eventHandler.UpdateEvent)
				adminEvents.DELETE("/:id", eventHandler.DeleteEvent)
			}
		}

		// Protected order routes
		orders := api.Group("/orders")
		orders.Use(middleware.AuthMiddleware())
		{
			orders.GET("", orderHandler.GetUserOrders)
			orders.GET("/:id", orderHandler.GetOrderByID)
			orders.POST("/:id/pay", orderHandler.ProcessPayment)
			orders.POST("/:id/cancel", orderHandler.CancelOrder)
		}
	}

	// ==========================================================
	// Step 8: Start the HTTP Server
	// ==========================================================
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Eventix API server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
