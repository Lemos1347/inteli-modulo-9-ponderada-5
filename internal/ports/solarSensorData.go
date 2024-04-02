package ports

import "github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"

type SolarSensorDataPort interface {
	CreateNewData(*entity.SolarSensorData)
}
