package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/domain/entity"
	"github.com/Lemos1347/inteli-modulo-9-ponderada-5/internal/ports"
	"github.com/google/uuid"
)

type SolarSensorRepo struct {
	solaSensorAdapter ports.SolarSensorPort
	mqttadapter       ports.MQTTPort
	names             []int
	coords            *map[float64]bool
	mutexMap          *sync.RWMutex
	mutexArr          *sync.Mutex
}

func (s *SolarSensorRepo) CreateNewSolarSensor() *entity.SolarSensor {
	x, y := s.generateCoords()
	solarSensor := entity.SolarSensor{
		Id:      uuid.New().String(),
		Name:    *s.generateName(),
		CoordsX: *x,
		CoordsY: *y,
	}

	s.solaSensorAdapter.CreateNewSensor(&solarSensor)

	return &solarSensor
}

func (s *SolarSensorRepo) generateName() *string {
	s.mutexArr.Lock()
	defer s.mutexArr.Unlock()

	size := len(s.names)

	if size == 0 {
		number := 1
		name := fmt.Sprintf("Solar_sensor %d\n", number)

		s.names = []int{
			number,
		}

		return &name
	}

	newNumber := s.names[size-1] + 1

	s.names = append(s.names, newNumber)

	newName := fmt.Sprintf("Solar_sensor %d\n", newNumber)

	return &newName

}

func (s *SolarSensorRepo) generateCoords() (*float64, *float64) {
	baseX := -46.633308
	baseY := -23.550520

	maxOffsetX, minOffsetX := 0.3, -0.2
	maxOffsetY, minOffsetY := 0.1, -0.2

	defer s.mutexMap.Unlock()

	for {
		offsetX := minOffsetX + (rand.Float64() * (maxOffsetX - minOffsetX))
		offsetY := minOffsetY + (rand.Float64() * (maxOffsetY - minOffsetY))

		coordsX := baseX + offsetX
		coordsY := baseY + offsetY

		s.mutexMap.RLock()
		if (*s.coords)[coordsX+coordsY] {
			s.mutexMap.RUnlock()
			continue
		} else {
			s.mutexMap.RUnlock()
			s.mutexMap.Lock()
			(*s.coords)[coordsX+coordsY] = true
			return &coordsX, &coordsY
		}
	}
}

func (s *SolarSensorRepo) ActivateExistingSensors() (*[]entity.SolarSensor, error) {
	existingSensors := s.solaSensorAdapter.GetAllSensors()

	if existingSensors == nil || len(*existingSensors) == 0 {
		return existingSensors, errors.New("No sensors in db available to activation")
	}

	wg := sync.WaitGroup{}

	for _, sensor := range *existingSensors {
		wg.Add(2)
		go s.setExistingCoords(sensor.CoordsX+sensor.CoordsY, &wg)
		go s.setExistingName(&sensor.Name, &wg)
	}

	wg.Wait()

	return existingSensors, nil
}

func (s *SolarSensorRepo) setExistingCoords(sumCoords float64, wgs ...*sync.WaitGroup) {
	if wgs != nil || len(wgs) > 0 {
		for _, wg := range wgs {
			defer wg.Done()
		}
	}

	s.mutexMap.Lock()
	defer s.mutexMap.Unlock()

	// TODO: dont repeat coords
	// if _, ok := (*s.coords)[sumCoords]; ok {
	// 	panic("Trying to activate coords that already exists\n")
	// }

	(*s.coords)[sumCoords] = true
}

func (s *SolarSensorRepo) setExistingName(name *string, wgs ...*sync.WaitGroup) {
	if wgs != nil || len(wgs) > 0 {
		for _, wg := range wgs {
			defer wg.Done()
		}
	}

	arrNames := strings.Split(strings.TrimSpace(*name), " ")

	if arrNames == nil || len(arrNames) == 0 {
		msgError := fmt.Sprintf("Couldn't split name: %s\n", *name)
		panic(msgError)
	}

	value, err := strconv.Atoi(arrNames[1])

	if err != nil {
		// log.Printf("Trying to get: %#v\n", arrNames)
		msgError := fmt.Sprintf("Last char of %s isn't a number\n%#v", *name, arrNames)
		panic(msgError)
	}

	s.mutexArr.Lock()
	defer s.mutexArr.Unlock()

	s.names = append(s.names, value)
}

func (s *SolarSensorRepo) emmitData(sensorId string, data string) {
	msg := entity.SolarSensorData{
		SensorId:  sensorId,
		Data:      data,
		CreatedAt: time.Now(),
	}

	jsonMsg, err := json.Marshal(msg)

	if err != nil {
		msgError := fmt.Sprintf("Error when encoding: %#v\n", msg)
		panic(msgError)
	}

	s.mqttadapter.Publish("sensors/data", 1, false, jsonMsg)
}

func (s *SolarSensorRepo) emulateSensor(sensor entity.SolarSensor) {
	for {
		time.Sleep(time.Second * 1)
		data, err := sensor.GenerateReading(os.Getenv("CSV_PATH"))

		if err != nil {
			msgError := fmt.Sprintf("Error generate reading of %s due to: %s", sensor.Name, err.Error())
			panic(msgError)
		}

		s.emmitData(sensor.Id, data)
	}
}

func (s *SolarSensorRepo) EmulateSensors(sensors *[]entity.SolarSensor) {
	for _, sensor := range *sensors {
		go s.emulateSensor(sensor)
	}

	select {}
}

func NewSolarSensorRepo(solarSensorAdapter ports.SolarSensorPort, mqttadapter ports.MQTTPort) *SolarSensorRepo {
	coords := make(map[float64]bool)
	mutexMap := sync.RWMutex{}
	mutexArr := sync.Mutex{}

	return &SolarSensorRepo{
		solaSensorAdapter: solarSensorAdapter,
		mqttadapter:       mqttadapter,
		names:             []int{},
		coords:            &coords,
		mutexMap:          &mutexMap,
		mutexArr:          &mutexArr,
	}
}
