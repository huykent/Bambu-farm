package domain

import (
	"time"

	"gorm.io/gorm"
)

type PrintJob struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrganizationID uint           `gorm:"index" json:"organization_id"`
	PrinterID      uint           `gorm:"index" json:"printer_id"`
	UserID         uint           `gorm:"index" json:"user_id"` // User who submitted the job
	Status         string         `gorm:"not null;default:'pending'" json:"status"` // pending, printing, paused, completed, failed, cancelled
	FileURI        string         `gorm:"not null" json:"file_uri"`
	FileName       string         `gorm:"not null" json:"file_name"`
	Progress       int            `json:"progress"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

type PrintHistory struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	PrintJobID  uint      `gorm:"index" json:"print_job_id"`
	Status      string    `json:"status"` // Status transition
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}
