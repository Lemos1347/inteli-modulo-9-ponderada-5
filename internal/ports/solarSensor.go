package ports

import "github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"

type SolarSensorPort interface {
	GetAllSensors() *[]entity.SolarSensor
	CreateNewSensor(solarSensor *entity.SolarSensor)
}
