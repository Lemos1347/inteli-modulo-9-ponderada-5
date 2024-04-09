package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/primary/queue"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/secondary/sensors"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/infra"
	"github.com/google/uuid"

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

func main() {
	dbConnection := infra.NewDBConnection()

	solarSensorDataAdapter := sensors.NewSolarSensorDataAdapter(dbConnection)

	kafkaConsumer := infra.GenerateKafkaConsumer()
	queueAdapter := queue.NewMessageHandler(kafkaConsumer, []string{"ponderada5"})

	callback := func(msg []byte) {
		log.Printf("Message received: %s", string(msg))

		var data entity.SolarSensorData

		err := json.Unmarshal(msg, &data)

		if err != nil {
			errMsg := fmt.Sprintf("Unable to decode json: %s\n", err.Error())
			panic(errMsg)
		}

		data.Id = uuid.New().String()

		solarSensorDataAdapter.CreateNewData(&data)
	}
	queueAdapter.ConsumeMessages(callback)

	// mqttAdapter.Subscribe("sensors/data", callback)

	select {}
}
