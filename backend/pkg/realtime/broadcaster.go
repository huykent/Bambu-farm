package realtime

import (
	"encoding/json"
)

type EventType string

const (
	PrinterStatusUpdate EventType = "printer_status_update"
	JobProgress         EventType = "job_progress"
	TemperatureUpdate   EventType = "temperature_update"
	Alerts              EventType = "alerts"
)

type Event struct {
	Type EventType   `json:"type"`
	Data interface{} `json:"data"`
}

type Broadcaster struct {
	manager *Manager
}

func NewBroadcaster(manager *Manager) *Broadcaster {
	return &Broadcaster{manager: manager}
}

func (b *Broadcaster) Publish(eventType EventType, data interface{}) {
	evt := Event{
		Type: eventType,
		Data: data,
	}

	payload, err := json.Marshal(evt)
	if err == nil {
		b.manager.broadcast <- payload
	}
}
