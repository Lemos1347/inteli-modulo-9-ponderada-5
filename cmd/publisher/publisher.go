package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/primary/mqtt"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/adapters/secondary/sensors"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/infra"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/repository"

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

func callback(_ MQTT.Client, msg MQTT.Message) {
	log.Printf("Message received: %s", string(msg.Payload()))
}

func main() {
	if os.Getenv("CSV_PATH") == "" {
		fmt.Println("\033[31mMissing csv path\033[0m")
		os.Exit(1)
	}

	dbConnection := infra.NewDBConnection()

	mqttClient := infra.GenerateMQTTClient()
	mqttAdapter := mqtt.NewMQTTAdapter(mqttClient)

	solarSensorAdapter := sensors.NewSolarSensorAdapter(dbConnection)

	solarSensorRepo := repository.NewSolarSensorRepo(solarSensorAdapter, mqttAdapter)

	sensorsArr := &[]entity.SolarSensor{}

	if existingSensors, err := solarSensorRepo.ActivateExistingSensors(); err != nil {
		errMsg := fmt.Sprintf("Error activating existing sensors: %s\n", err.Error())
		log.Print(errMsg)
	} else {
		sensorsArr = existingSensors
	}

	switch {
	case len(os.Args) > 2:
		limit, err := strconv.Atoi(os.Args[2])
		if err != nil {
			panic("Sensors amount must be an integer")
		}

		for i := 0; i < limit; i++ {
			sensor := solarSensorRepo.CreateNewSolarSensor()
			*sensorsArr = append(*sensorsArr, *sensor)
		}
		solarSensorRepo.EmulateSensors(sensorsArr)

	case len(os.Args) <= 2 && len(*sensorsArr) == 0:
		sensor := solarSensorRepo.CreateNewSolarSensor()
		*sensorsArr = append(*sensorsArr, *sensor)
		solarSensorRepo.EmulateSensors(sensorsArr)

	default:
		solarSensorRepo.EmulateSensors(sensorsArr)

	}

	// if len(os.Args) > 2 {
	// 	limit, err := strconv.Atoi(os.Args[2])
	// 	if err != nil {
	// 		panic("Sensors amount must be an integer")
	// 	}
	//
	// 	for i := 0; i < limit; i++ {
	// 		sensor := solarSensorRepo.CreateNewSolarSensor()
	// 		*sensorsArr = append(*sensorsArr, *sensor)
	// 	}
	// } else {
	//    sensor := solarSensorRepo.CreateNewSolarSensor()
	//    *sensorsArr = append(*sensorsArr, sensor)
	//
	// sensorsArr := []entity.SolarSensor{
	// 	*sensor,
	// }
	//
	//  }
	//
	// // sensor := solarSensorRepo.CreateNewSolarSensor()
	// //
	// // sensorsArr := []entity.SolarSensor{
	// // 	*sensor,
	// // }
	//
	// solarSensorRepo.EmulateSensors(sensorsArr)
}
