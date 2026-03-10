package domain

import (
	"time"

	"gorm.io/gorm"
)

type Alert struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	OrganizationID uint           `gorm:"index" json:"organization_id"`
	PrinterID      uint           `gorm:"index" json:"printer_id"`
	Level          string         `json:"level"`   // info, warning, critical
	Type           string         `json:"type"`    // print_failure, temp_anomaly, offline
	Message        string         `json:"message"`
	Resolved       bool           `json:"resolved" gorm:"default:false"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
