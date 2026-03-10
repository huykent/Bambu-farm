package domain

import (
	"time"

	"gorm.io/gorm"
)

type Printer struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	OrganizationID  uint           `gorm:"index" json:"organization_id"`
	Organization    Organization   `json:"organization,omitempty"`
	PrinterID       string         `gorm:"uniqueIndex;not null" json:"printer_id"`
	Name            string         `gorm:"not null" json:"name"`
	IPAddress       string         `gorm:"not null" json:"ip_address"`
	AccessToken     string         `gorm:"not null" json:"-"` // encrypted in real app
	Model           string         `gorm:"not null" json:"model"`
	Status          string         `gorm:"not null;default:'offline'" json:"status"`
	FirmwareVersion string         `json:"firmware_version"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

type PrinterStatus struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	PrinterID      uint      `gorm:"index" json:"printer_id"`
	NozzleTemp     float64   `json:"nozzle_temp"`
	BedTemp        float64   `json:"bed_temp"`
	PrintProgress  int       `json:"print_progress"`
	Layer          int       `json:"layer"`
	RemainingTime  int       `json:"remaining_time"`
	JobState       string    `json:"job_state"`
	FilamentType   string    `json:"filament_type"`
	FanSpeed       int       `json:"fan_speed"`
	RecordedAt     time.Time `gorm:"autoCreateTime" json:"recorded_at"`
}

type PrinterLog struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	PrinterID  uint      `gorm:"index" json:"printer_id"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	ErrorCode  string    `json:"error_code,omitempty"`
	RecordedAt time.Time `gorm:"autoCreateTime" json:"recorded_at"`
}

type PrinterMetric struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	PrinterID uint      `gorm:"index" json:"printer_id"`
	MetricKey string    `gorm:"index" json:"metric_key"` // e.g. "nozzle_temp", "progress"
	Value     float64   `json:"value"`
	Timestamp time.Time `gorm:"autoCreateTime" json:"timestamp"`
}

