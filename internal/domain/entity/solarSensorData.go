package entity

import "time"

type SolarSensorData struct {
	Id        string `json:"id"`
	SensorId  string `json:"sensor_id"`
	Data      string `json:"data"`
	CreatedAt time.Time `json:"createdAt"`
}
