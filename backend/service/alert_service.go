package service

import (
	"bambu-farm/domain"
	"bambu-farm/pkg/alerting"
	"os"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AlertService struct {
	logger *zap.SugaredLogger
	db     *gorm.DB
}

func NewAlertService(logger *zap.SugaredLogger, db *gorm.DB) *AlertService {
	return &AlertService{logger: logger, db: db}
}

func (s *AlertService) HandleAlert(orgID, printerID uint, level, alertType, message string) error {
	// Save to DB
	alert := domain.Alert{
		OrganizationID: orgID,
		PrinterID:      printerID,
		Level:          level,
		Type:           alertType,
		Message:        message,
	}

	if err := s.db.Create(&alert).Error; err != nil {
		s.logger.Errorf("Failed to save alert to DB: %v", err)
		return err
	}

	// Read organization alert settings (Stubbed with env vars for prototype)
	tgToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	tgChatID := os.Getenv("TELEGRAM_CHAT_ID")
	alertEmail := os.Getenv("ALERT_EMAIL")

	if tgToken != "" && tgChatID != "" {
		alerting.SendTelegramAlert(s.logger, tgToken, tgChatID, message)
	}

	if alertEmail != "" {
		alerting.SendEmailAlert(s.logger, alertEmail, "BambuFarm Alert - "+level, message)
	}

	return nil
}
