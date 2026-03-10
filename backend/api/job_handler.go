package api

import (
	"bambu-farm/pkg/auth"
	"bambu-farm/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type JobHandler struct {
	jobService *service.JobService
}

func NewJobHandler(jobService *service.JobService) *JobHandler {
	return &JobHandler{jobService: jobService}
}

func (h *JobHandler) RegisterRoutes(router *gin.Engine) {
	jobGroup := router.Group("/jobs")
	jobGroup.Use(auth.JWTMiddleware())

	jobGroup.GET("", h.ListJobs)
	jobGroup.POST("/submit", h.SubmitJob)
	
	// Actions
	jobGroup.POST("/:id/pause", h.ActionJob("pause"))
	jobGroup.POST("/:id/resume", h.ActionJob("resume"))
	jobGroup.POST("/:id/cancel", h.ActionJob("cancel"))
}

type SubmitJobRequest struct {
	PrinterID uint   `json:"printer_id" binding:"required"`
	FileURI   string `json:"file_uri" binding:"required"`
	FileName  string `json:"file_name" binding:"required"`
}

func (h *JobHandler) SubmitJob(c *gin.Context) {
	orgID := getUintFromCtx(c, "organizationID")
	userID := getUintFromCtx(c, "userID")

	var req SubmitJobRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.jobService.SubmitJob(orgID, userID, req.PrinterID, req.FileURI, req.FileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Job submitted", "job": job})
}

func (h *JobHandler) ListJobs(c *gin.Context) {
	orgID := getUintFromCtx(c, "organizationID")
	
	jobs, err := h.jobService.ListJobs(orgID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"jobs": jobs})
}

func (h *JobHandler) ActionJob(action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		orgID := getUintFromCtx(c, "organizationID")
		idStr := c.Param("id")
		
		id, err := strconv.ParseUint(idStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid job id"})
			return
		}

		if err := h.jobService.UpdateJobStatus(orgID, uint(id), action); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Job action successful", "action": action})
	}
}

func getUintFromCtx(c *gin.Context, key string) uint {
	val, exists := c.Get(key)
	if !exists {
		return 0
	}
	return val.(uint)
}
