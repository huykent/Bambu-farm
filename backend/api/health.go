package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes sets up all the API endpoints
func RegisterRoutes(router *gin.Engine) {
	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to BambuLab Print Farm Manager API",
		})
	})

	// Health check endpoint (Required by Module 1)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
}
