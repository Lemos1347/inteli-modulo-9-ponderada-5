package sensors

import (
	"fmt"

	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"
	"gorm.io/gorm"
)

type SolarSensorAdapter struct {
	db *gorm.DB
}

func (s *SolarSensorAdapter) GetAllSensors() *[]entity.SolarSensor {
	var solarSensors []entity.SolarSensor
	result := s.db.Select("id", "name").Find(&solarSensors)

	if result.Error != nil {
		panic(fmt.Sprintf("Unable to get all records from SolarSensor table due to: %s", result.Error.Error()))
	}

	return &solarSensors
}

func (s *SolarSensorAdapter) CreateNewSensor(solarSensor *entity.SolarSensor) {
	result := s.db.Create(solarSensor)

	if result.Error != nil {
		panic(fmt.Sprintf("Unable to create solaSensor %s due to: %s", solarSensor.Name, result.Error.Error()))
	}
}

func NewSolarSensorAdapter(dbConnection *gorm.DB) *SolarSensorAdapter {
	return &SolarSensorAdapter{
		db: dbConnection,
	}
}
