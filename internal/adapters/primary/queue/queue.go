package queue

import (
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type callbackFunc func([]byte)

type MessageHandler struct {
	consumer *kafka.Consumer
}

func (s *MessageHandler) ConsumeMessages(callback callbackFunc, oneTime ...bool) {
	log.Print("Listing for messages!")
	defer s.consumer.Close()
	for {
		// consumes messages from the subscribed topic and prints them to the console
		e := s.consumer.Poll(1000)
		switch ev := e.(type) {
		case *kafka.Message:
			// application-specific processing
			fmt.Printf("Consumed event from topic %s: key = %-10s value = %s\n",
				*ev.TopicPartition.Topic, string(ev.Key), string(ev.Value))

			callback(ev.Value)
		case kafka.Error:
			fmt.Fprintf(os.Stderr, "%% Error: %v\n", ev)
			break
		}

		if len(oneTime) > 0 {
			break
		}

	}
}

func NewMessageHandler(consumer *kafka.Consumer, topics []string) *MessageHandler {
	err := consumer.SubscribeTopics(topics, nil)

	if err != nil {
		log.Fatalf("Error trying to subscribe into topics: %s\n", err.Error())
	}

	return &MessageHandler{
		consumer: consumer,
	}
}
