package api

import (
	"bambu-farm/pkg/auth"
	"bambu-farm/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PrinterHandler struct {
	printerService *service.PrinterService
}

func NewPrinterHandler(printerService *service.PrinterService) *PrinterHandler {
	return &PrinterHandler{printerService: printerService}
}

func (h *PrinterHandler) RegisterRoutes(router *gin.Engine) {
	printerGroup := router.Group("/printers")
	// Require authentication for all printer routes
	printerGroup.Use(auth.JWTMiddleware())

	printerGroup.GET("", h.ListPrinters)
	printerGroup.POST("", h.AddPrinter)
	printerGroup.GET("/:id", h.GetPrinter)
	printerGroup.DELETE("/:id", h.DeletePrinter)
}

type AddPrinterRequest struct {
	PrinterID   string `json:"printer_id" binding:"required"`
	Name        string `json:"name" binding:"required"`
	IPAddress   string `json:"ip_address" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	Model       string `json:"model" binding:"required"`
}

func (h *PrinterHandler) AddPrinter(c *gin.Context) {
	orgIDVal, exists := c.Get("organizationID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "organization missing from context"})
		return
	}
	orgID := orgIDVal.(uint)

	var req AddPrinterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	printer, err := h.printerService.AddPrinter(orgID, req.PrinterID, req.Name, req.IPAddress, req.AccessToken, req.Model)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Printer added successfully", "printer": printer})
}

func (h *PrinterHandler) ListPrinters(c *gin.Context) {
	orgIDVal, _ := c.Get("organizationID")
	orgID := orgIDVal.(uint)

	printers, err := h.printerService.ListPrinters(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"printers": printers})
}

func (h *PrinterHandler) GetPrinter(c *gin.Context) {
	orgIDVal, _ := c.Get("organizationID")
	orgID := orgIDVal.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid printer id"})
		return
	}

	printer, err := h.printerService.GetPrinter(uint(id), orgID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"printer": printer})
}

func (h *PrinterHandler) DeletePrinter(c *gin.Context) {
	orgIDVal, _ := c.Get("organizationID")
	orgID := orgIDVal.(uint)

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid printer id"})
		return
	}

	if err := h.printerService.DeletePrinter(uint(id), orgID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Printer deleted successfully"})
}
