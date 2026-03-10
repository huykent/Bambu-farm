package telemetry

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"bambu-farm/domain"
	"bambu-farm/service"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Collector struct {
	logger         *zap.SugaredLogger
	db             *gorm.DB
	printerService *service.PrinterService
	clients        map[uint]*MQTTClient
	mu             sync.RWMutex
}

func NewCollector(logger *zap.SugaredLogger, db *gorm.DB, printerService *service.PrinterService) *Collector {
	return &Collector{
		logger:         logger,
		db:             db,
		printerService: printerService,
		clients:        make(map[uint]*MQTTClient),
	}
}

// Start initializes connections to all known printers
func (c *Collector) Start(ctx context.Context) {
	c.logger.Info("Starting Telemetry Collector...")

	// In a real application, we would query active organizations, then printers.
	// For this prototype, we'll assume we have a way to get all configured printers
	// Alternatively, we run a routine that periodically checks DB for new printers.
	
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		defer ticker.Stop()
		
		c.syncPrinters()

		for {
			select {
			case <-ctx.Done():
				c.stopAll()
				return
			case <-ticker.C:
				c.syncPrinters()
			}
		}
	}()
}

func (c *Collector) syncPrinters() {
	// Dummy query: fetch all printers. In real usage, be organization-aware.
	var printers []domain.Printer
	if err := c.db.Find(&printers).Error; err != nil {
		c.logger.Errorf("Failed to fetch printers for telemetry: %v", err)
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	for _, p := range printers {
		if _, exists := c.clients[p.ID]; !exists && p.IPAddress != "" && p.AccessToken != "" {
			client, err := NewMQTTClient(c.logger, p.IPAddress, p.PrinterID, p.AccessToken)
			if err != nil {
				c.logger.Errorf("Failed to connect to printer %d (%s): %v", p.ID, p.IPAddress, err)
				continue
			}

			topic := fmt.Sprintf("device/%s/report", p.PrinterID)
			err = client.Subscribe(topic, c.createMessageHandler(p.ID))
			if err != nil {
				c.logger.Errorf("Failed to subscribe to topic %s: %v", topic, err)
				client.Disconnect()
				continue
			}

			c.clients[p.ID] = client
			c.logger.Infof("Telemetry active for printer %d", p.ID)
		}
	}
}

func (c *Collector) createMessageHandler(printerID uint) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		// Parse BambuLab JSON payload
		var payload map[string]interface{}
		if err := json.Unmarshal(msg.Payload(), &payload); err != nil {
			c.logger.Errorf("Failed to unmarshal telemetry for printer %d: %v", printerID, err)
			return
		}

		// Example parsing logic (Bambu structure is deeply nested inside "print")
		if printData, ok := payload["print"].(map[string]interface{}); ok {
			
			// Extract Nozzle Temp
			if nozzleTarget, ok := printData["nozzle_temper"].(float64); ok {
				c.saveMetric(printerID, "nozzle_temp", nozzleTarget)
			}
			
			// Extract Bed Temp
			if bedTarget, ok := printData["bed_temper"].(float64); ok {
				c.saveMetric(printerID, "bed_temp", bedTarget)
			}
			
			// Extract Progress
			if progress, ok := printData["mc_percent"].(float64); ok {
				c.saveMetric(printerID, "progress", progress)
			}
		}
	}
}

func (c *Collector) saveMetric(printerID uint, key string, value float64) {
	metric := domain.PrinterMetric{
		PrinterID: printerID,
		MetricKey: key,
		Value:     value,
	}
	if err := c.db.Create(&metric).Error; err != nil {
		c.logger.Errorf("Failed to save metric %s for printer %d: %v", key, printerID, err)
	}
}

func (c *Collector) stopAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for _, client := range c.clients {
		client.Disconnect()
	}
	c.logger.Info("Telemetry Collector stopped")
}
