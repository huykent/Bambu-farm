package telemetry

import (
	"crypto/tls"
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type MQTTClient struct {
	logger *zap.SugaredLogger
	client mqtt.Client
}

func NewMQTTClient(logger *zap.SugaredLogger, brokerIP, serialNumber, accessCode string) (*MQTTClient, error) {
	// BambuLab printers use MQTTS on port 8883
	brokerURI := fmt.Sprintf("tls://%s:8883", brokerIP)

	opts := mqtt.NewClientOptions()
	opts.AddBroker(brokerURI)
	opts.SetClientID(fmt.Sprintf("bambu-manager-%d", time.Now().UnixNano()))
	opts.SetUsername("bblp")
	opts.SetPassword(accessCode)

	// Insecure skip verify because Bambus use self-signed certs
	opts.SetTLSConfig(&tls.Config{InsecureSkipVerify: true})
	opts.SetAutoReconnect(true)
	opts.SetMaxReconnectInterval(10 * time.Second)

	client := mqtt.NewClient(opts)
	token := client.Connect()
	if token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	logger.Infof("Connected to MQTT broker at %s", brokerIP)

	return &MQTTClient{
		logger: logger,
		client: client,
	}, nil
}

func (m *MQTTClient) Subscribe(topic string, handler mqtt.MessageHandler) error {
	token := m.client.Subscribe(topic, 1, handler)
	if token.Wait() && token.Error() != nil {
		return token.Error()
	}
	m.logger.Infof("Subscribed to MQTT topic: %s", topic)
	return nil
}

func (m *MQTTClient) Disconnect() {
	m.client.Disconnect(250)
}
