package kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"message/internal/storage"
)

var Producer *kafka.Producer

func StartKafka() {
	var err error
	Producer, err = kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": fmt.Sprintf("localhost:%s", storage.Env.KafkaPort),
	})
	if err != nil {
		log.Fatal(err)
	}

	defer Producer.Close()
}
