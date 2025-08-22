package kafka

import (
	"fmt"
	kafkago "github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"message/internal/storage"
	"results/errs"
)

var Producer *kafkago.Producer

func StartKafka() error {
	var err error
	Producer, err = kafkago.NewProducer(&kafkago.ConfigMap{
		"bootstrap.servers":  fmt.Sprintf("%s:%s", storage.Env.KafkaHost, storage.Env.KafkaPort),
		"message.timeout.ms": 5000,
		"acks":               "all",
	})
	if err != nil {
		return err
	}

	go func() {
		for e := range Producer.Events() {
			switch ev := e.(type) {
			case *kafkago.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("%s: %v", errs.KafkaDeliverFailed, ev.TopicPartition.Error)
				}
			case kafkago.Error:
				log.Println(ev)
			}
		}
	}()

	return nil
}
