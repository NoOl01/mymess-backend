package kafka_consumer

import (
	"chat_kafka_broker/internal/grpc_client"
	"context"
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"proto/scyllapb"
	"results/errs"
	"results/succ"
	"time"
)

type KafkaJson struct {
	Message string    `json:"message"`
	UserId  int64     `json:"user_id"`
	ChatId  string    `json:"chat_id"`
	Time    time.Time `json:"time"`
}

type Consumer struct {
	KafkaConsumer *kafka.Consumer
	topic         string
	batchSize     int
	flushInterval time.Duration
	Client        scyllapb.ScyllaServiceClient
}

func NewConsumer(brokers, groupId, topic string, batchSize int, flushInterval time.Duration, scClient scyllapb.ScyllaServiceClient) (*Consumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  brokers,
		"group.id":           groupId,
		"auto.offset.reset":  "latest",
		"enable.auto.commit": false,
	})
	if err != nil {
		return nil, err
	}

	err = c.SubscribeTopics([]string{topic}, nil)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		KafkaConsumer: c,
		topic:         topic,
		batchSize:     batchSize,
		flushInterval: flushInterval,
		Client:        scClient,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	var (
		batch         []KafkaJson
		messageBuffer []*kafka.Message
		lastFlush     = time.Now()
	)

	for {
		select {
		case <-ctx.Done():
			log.Printf("%s:", errs.KafkaConsumerContextCanceled)
			c.flush(batch, messageBuffer)
			_ = c.KafkaConsumer.Close()
			return
		default:
			ev, err := c.KafkaConsumer.ReadMessage(5 * time.Second)
			if err != nil {
				log.Println(err)
				continue
			}

			var msg KafkaJson
			if err := json.Unmarshal(ev.Value, &msg); err != nil {
				log.Println(err)
				continue
			}

			batch = append(batch, msg)
			messageBuffer = append(messageBuffer, ev)

			if len(batch) >= c.batchSize || time.Since(lastFlush) > c.flushInterval {
				c.flush(batch, messageBuffer)
				batch = nil
				messageBuffer = nil
				lastFlush = time.Now()
			}
		}
	}
}

func (c *Consumer) flush(batch []KafkaJson, messages []*kafka.Message) {
	if len(batch) == 0 {
		return
	}

	var grpcMessages []*scyllapb.ChatMessage

	for _, message := range batch {
		grpcMessages = append(grpcMessages, &scyllapb.ChatMessage{
			Message: message.Message,
			UserId:  message.UserId,
			ChatId:  message.ChatId,
			Time:    timestamppb.New(message.Time),
		})
	}

	resp, err := grpc_client.UploadMessages(c.Client, grpcMessages)
	if err != nil || resp.Result != succ.Ok {
		log.Printf("Upload failed: %v", err)
		return
	}

	lastMsg := messages[len(batch)-1]
	_, err = c.KafkaConsumer.CommitOffsets([]kafka.TopicPartition{
		{
			Topic:     lastMsg.TopicPartition.Topic,
			Partition: lastMsg.TopicPartition.Partition,
			Offset:    lastMsg.TopicPartition.Offset + 1,
		},
	})
	if err != nil {
		log.Println(err)
	}
}
