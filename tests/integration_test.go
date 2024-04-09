package tests

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/primary/mqtt"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/primary/queue"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/secondary/sensors"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/infra"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/repository"
	"github.com/joho/godotenv"
)

func TestKafkaConnection(t *testing.T) {
	// TODO: pass .env vars
	err := godotenv.Load("../configs/.env.test")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("ENV variables loaded!")
	consumer := infra.GenerateKafkaConsumer()
	queue.NewMessageHandler(consumer, []string{"ponderada5_test"})

	t.Log("Kafka consumer connected!")
	return
}

func TestMQTTConnection(t *testing.T) {
	// TODO: pass .env vars
	godotenv.Load("../configs/.env.test")
	log.Println("ENV variables loaded!")

	mqttClient := infra.GenerateMQTTClient()
	mqtt.NewMQTTAdapter(mqttClient)

	t.Log("MQTTClient connected!")
	return
}

func TestIntegration(t *testing.T) {
	godotenv.Load("../configs/.env.test")
	log.Println("ENV variables loaded!")
	var nothing *sensors.SolarSensorAdapter

	mqttClient := infra.GenerateMQTTClient()
	mqttAdapter := mqtt.NewMQTTAdapter(mqttClient)

	solarSensor := repository.NewSolarSensorRepo(nothing, mqttAdapter)

	msgEmitted := solarSensor.EmmitData("sensor-01", "32", "ponderada5/test")
	t.Logf("Msg emmited: %#v\n", msgEmitted)

	consumer := infra.GenerateKafkaConsumer()
	queueMsgsHandler := queue.NewMessageHandler(consumer, []string{"ponderada5_test"})

	checkMsgs := func(msgFromQueue []byte) {
		var msgReceived entity.SolarSensorData

		json.Unmarshal(msgFromQueue, &msgReceived)

		if msgReceived.CreatedAt != msgEmitted.CreatedAt || msgReceived.SensorId != msgEmitted.SensorId || msgReceived.Data != msgEmitted.Data {
			t.Fatal("Different messages!")
		}

	}
	queueMsgsHandler.ConsumeMessages(checkMsgs, true)

	t.Log("Messages are the same!")
}
