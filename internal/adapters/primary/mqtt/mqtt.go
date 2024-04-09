package mqtt

import (
	"log"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

type MQTTAdapter struct {
	client MQTT.Client
}

func (s *MQTTAdapter) Publish(topic string, qos byte, retained bool, data interface{}) {
	token := s.client.Publish(topic, qos, retained, data)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Error trying to publish message in MQTT topic: %s\n", token.Error().Error())
	}
}

func (s *MQTTAdapter) Subscribe(topic string, callback MQTT.MessageHandler) {
	token := s.client.Subscribe(topic, 1, callback)
	token.Wait()
	if token.Error() != nil {
		log.Fatalf("Error trying to subscribe MQTT topic: %s\n", token.Error().Error())
	}
	log.Println("CLient subscribed sucessifuly!")
}

func NewMQTTAdapter(client MQTT.Client) *MQTTAdapter {
	return &MQTTAdapter{
		client: client,
	}
}
