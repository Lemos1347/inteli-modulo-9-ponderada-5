package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/primary/mqtt"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/secondary/sensors"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/infra"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/ports"
	"github.com/google/uuid"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
)

func init() {
	log.Print("Loading .env variables...")

	envPath := "./.env"

	if len(os.Args) > 1 {
		envPath = os.Args[1]
	}

	godotenv.Load(envPath)

	log.Println("ENV variables loaded!")
}

func createCallback(solarSensorDataAdapter ports.SolarSensorDataPort) MQTT.MessageHandler {

	return func(_ MQTT.Client, msg MQTT.Message) {
		log.Printf("Message received: %s", string(msg.Payload()))

		var data entity.SolarSensorData

		err := json.Unmarshal(msg.Payload(), &data)

		if err != nil {
			errMsg := fmt.Sprintf("Unable to decode json: %s\n", err.Error())
			panic(errMsg)
		}

		data.Id = uuid.New().String()

		solarSensorDataAdapter.CreateNewData(&data)
	}
}

func main() {
	dbConnection := infra.NewDBConnection()

	solarSensorDataAdapter := sensors.NewSolarSensorDataAdapter(dbConnection)

	mqttClient := infra.GenerateMQTTClient()
	mqttAdapter := mqtt.NewMQTTAdapter(mqttClient)

	callback := createCallback(solarSensorDataAdapter)

	mqttAdapter.Subscribe("sensors/data", callback)

	select {}
}
