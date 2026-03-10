package alerting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"os"

	"go.uber.org/zap"
)

func SendTelegramAlert(logger *zap.SugaredLogger, botToken, chatID, message string) {
	if botToken == "" || chatID == "" {
		return
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	
	payload := map[string]string{
		"chat_id": chatID,
		"text":    "🚨 BambuFarm Alert 🚨\n" + message,
	}
	
	body, _ := json.Marshal(payload)
	
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		logger.Errorf("Failed to send Telegram alert: %v", err)
		return
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("Telegram API returned non-200 status: %d", resp.StatusCode)
	}
}

func SendEmailAlert(logger *zap.SugaredLogger, toEmail, subject, message string) {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPass := os.Getenv("SMTP_PASS")

	if smtpHost == "" || toEmail == "" {
		return
	}

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	msg := []byte("To: " + toEmail + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		message + "\r\n")

	addr := fmt.Sprintf("%s:%s", smtpHost, smtpPort)
	err := smtp.SendMail(addr, auth, smtpUser, []string{toEmail}, msg)
	if err != nil {
		logger.Errorf("Failed to send Email alert: %v", err)
	}
}
