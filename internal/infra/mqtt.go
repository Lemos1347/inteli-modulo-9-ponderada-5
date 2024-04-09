package infra

import (
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

func GenerateMQTTClient() MQTT.Client {

	opts := MQTT.NewClientOptions().AddBroker(os.Getenv("BROKER_URL"))
	opts.SetClientID(uuid.New().String())
  opts.SetUsername(os.Getenv("BROKER_USER"))
	opts.SetPassword(os.Getenv("BROKER_PASSWORD"))

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client
}
