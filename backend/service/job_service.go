package service

import (
	"bambu-farm/domain"
	"bambu-farm/pkg/queue"
	"errors"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type JobService struct {
	db  *gorm.DB
	rdb *redis.Client
}

func NewJobService(db *gorm.DB, rdb *redis.Client) *JobService {
	return &JobService{db: db, rdb: rdb}
}

func (s *JobService) SubmitJob(orgID, userID, printerID uint, fileURI, fileName string) (*domain.PrintJob, error) {
	job := &domain.PrintJob{
		OrganizationID: orgID,
		PrinterID:      printerID,
		UserID:         userID,
		Status:         "pending",
		FileURI:        fileURI,
		FileName:       fileName,
	}

	if err := s.db.Create(job).Error; err != nil {
		return nil, err
	}

	s.logHistory(job.ID, "Job Submitted")

	// Enqueue to Redis
	cmd := queue.PrinterCommand{
		JobID:     job.ID,
		PrinterID: printerID,
		Command:   "start",
		Payload:   fileURI,
	}
	if err := queue.EnqueueCommand(s.rdb, cmd); err != nil {
		return nil, errors.New("failed to enqueue job")
	}

	return job, nil
}

func (s *JobService) UpdateJobStatus(orgID, jobID uint, newStatus string) error {
	var job domain.PrintJob
	if err := s.db.Where("id = ? AND organization_id = ?", jobID, orgID).First(&job).Error; err != nil {
		return err
	}

	job.Status = newStatus
	if err := s.db.Save(&job).Error; err != nil {
		return err
	}

	s.logHistory(job.ID, "Status changed to: "+newStatus)

	cmd := queue.PrinterCommand{
		JobID:     job.ID,
		PrinterID: job.PrinterID,
		Command:   newStatus, // e.g. "pause", "resume", "cancel"
	}
	_ = queue.EnqueueCommand(s.rdb, cmd)

	return nil
}

func (s *JobService) ListJobs(orgID uint) ([]domain.PrintJob, error) {
	var jobs []domain.PrintJob
	err := s.db.Where("organization_id = ?", orgID).Order("created_at desc").Find(&jobs).Error
	return jobs, err
}

func (s *JobService) logHistory(jobID uint, notes string) {
	history := domain.PrintHistory{
		PrintJobID: jobID,
		Notes:      notes,
	}
	s.db.Create(&history)
}
