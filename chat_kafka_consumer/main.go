package main

import (
	"chat_kafka_broker/internal/grpc_client"
	kafkaconsumer "chat_kafka_broker/internal/kafka"
	"chat_kafka_broker/internal/storage"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"os/signal"
	"results/errs"
	"syscall"
	"time"
)

func main() {
	storage.LoadEnv()

	scClient, scConn, err := grpc_client.GrpcScyllaClientConnect()
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("%s: %v", errs.GrpcClientCloseFailed, err)
		}
	}(scConn)

	consumer, err := kafkaconsumer.NewConsumer(fmt.Sprintf("%s:%s", storage.Env.KafkaHost, storage.Env.KafkaPort), storage.Env.KafkaGroupId, storage.Env.KafkaTopic, 1000, 5*time.Second, scClient)
	if err != nil {
		log.Fatalf(err.Error())
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		stop()
		_ = consumer.KafkaConsumer.Close()
	}()

	consumer.Start(ctx)
}
