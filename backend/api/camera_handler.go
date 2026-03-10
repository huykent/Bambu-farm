package api

import (
	"bambu-farm/pkg/auth"
	"bambu-farm/pkg/camera"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CameraHandler struct {
	proxyService *camera.ProxyService
}

func NewCameraHandler(proxyService *camera.ProxyService) *CameraHandler {
	return &CameraHandler{proxyService: proxyService}
}

func (h *CameraHandler) RegisterRoutes(router *gin.Engine) {
	cameraGroup := router.Group("/printers/:id/camera")
	cameraGroup.Use(auth.JWTMiddleware())

	cameraGroup.GET("/stream", h.StreamCamera)
}

func (h *CameraHandler) StreamCamera(c *gin.Context) {
	orgIDVal, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "organization missing from context"})
		return
	}
	orgID := orgIDVal.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid printer id"})
		return
	}

	h.proxyService.StreamHandler(c, orgID, uint(id))
}
