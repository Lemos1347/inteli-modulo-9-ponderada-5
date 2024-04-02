package ports

import MQTT "github.com/eclipse/paho.mqtt.golang"

type MQTTPort interface {
	Publish(string, byte, bool, interface{})
	Subscribe(string, MQTT.MessageHandler)
}
