package infra

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func readConfluentConfig() kafka.ConfigMap {
	// reads the client configuration from client.properties
	// and returns it as a key-value map
	m := make(map[string]kafka.ConfigValue)

	file, err := os.Open(os.Getenv("PROPERTIES_FILE_PATH"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open file: %s", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, "#") && len(line) != 0 {
			kv := strings.Split(line, "=")
			parameter := strings.TrimSpace(kv[0])
			value := strings.TrimSpace(kv[1])
			m[parameter] = value
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Failed to read file: %s", err)
		os.Exit(1)
	}

	return m
}

func GenerateKafkaConsumer() *kafka.Consumer {
	conf := readConfluentConfig()
	conf["group.id"] = "go-group-1"
	conf["auto.offset.reset"] = "earliest"
	consumer, err := kafka.NewConsumer(&conf)

	if err != nil {
		log.Fatalf("Error generating kafka consumer: %s\n", err.Error())
	}

	return consumer
}
