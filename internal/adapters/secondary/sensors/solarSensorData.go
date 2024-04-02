package sensors

import (
	"fmt"

	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"
	"gorm.io/gorm"
)

type SolarSensorData struct {
	db *gorm.DB
}

func (s *SolarSensorData) CreateNewData(data *entity.SolarSensorData) {
	result := s.db.Create(*data)

	if result.Error != nil {
		panic(fmt.Sprintf("Unable to create data from sensor %s due to: %s", data.SensorId, result.Error.Error()))
	}
}

func NewSolarSensorDataAdapter(db *gorm.DB) *SolarSensorData {
	return &SolarSensorData{
		db: db,
	}
}
