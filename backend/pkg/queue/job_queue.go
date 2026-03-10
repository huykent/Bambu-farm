package queue

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type PrinterCommand struct {
	JobID     uint   `json:"job_id"`
	PrinterID uint   `json:"printer_id"`
	Command   string `json:"command"` // start, pause, resume, cancel
	Payload   string `json:"payload,omitempty"` // URL to gcode, etc.
}

const queueName = "printer_commands"

func EnqueueCommand(rdb *redis.Client, cmd PrinterCommand) error {
	ctx := context.Background()
	data, err := json.Marshal(cmd)
	if err != nil {
		return err
	}

	return rdb.LPush(ctx, queueName, data).Err()
}

// StartWorker listens for commands and processes them.
// In a full implementation, this would connect to the PrinterService/MQTT wrapper and execute them.
func StartWorker(logger *zap.SugaredLogger, rdb *redis.Client) {
	// Stubbing the background worker
	go func() {
		ctx := context.Background()
		for {
			// BRPop blocks until an item is available
			result, err := rdb.BRPop(ctx, 0, queueName).Result()
			if err != nil {
				logger.Errorf("Redis Queue error: %v", err)
				time.Sleep(2 * time.Second)
				continue
			}

			// result[0] is the queue name, result[1] is the data
			var cmd PrinterCommand
			if err := json.Unmarshal([]byte(result[1]), &cmd); err != nil {
				logger.Errorf("Failed to parse queue command: %v", err)
				continue
			}

			logger.Infof("Processing Job ID: %d, Command: %s for Printer ID: %d", cmd.JobID, cmd.Command, cmd.PrinterID)
			
			// TODO: Forward command to actual printer via MQTT connection manager
		}
	}()
}
