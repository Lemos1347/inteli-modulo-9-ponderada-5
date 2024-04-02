package sensors

import (
	"encoding/csv"
	"fmt"
	"math/rand"
	"os"
	"time"
)

// function to determine the period of the day
func getTimePeriod() string {
	horaAtual := time.Now().Hour()
	switch {
	case horaAtual >= 6 && horaAtual < 12:
		return "Manhã"
	case horaAtual >= 12 && horaAtual < 18:
		return "Tarde"
	default:
		return "Noite"
	}
}

// function to generate reading of a given sensor
func GenerateReading(pathCSV string) (string, error) {
	file, err := os.Open(pathCSV)
	if err != nil {
		return "", err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvData, err := csvReader.ReadAll()
	if err != nil {
		return "", err
	}

	dayPeriod := getTimePeriod()
	for i, record := range csvData[1:] { 
		period := record[0]
		if period == dayPeriod {
			return csvData[i+rand.Intn(10)][1], nil
		}
	}

	return "", fmt.Errorf("\033[31mNão foram encontrados mais dados para o período %s \033[0m", dayPeriod)
}
