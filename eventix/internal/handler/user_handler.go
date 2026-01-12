package handler

import (
	"net/http"

	"eventix/internal/middleware"

	"github.com/gin-gonic/gin"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	role := middleware.GetUserRole(c)

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile retrieved successfully",
		"user": gin.H{
			"id":   userID,
			"role": role,
		},
	})
}
