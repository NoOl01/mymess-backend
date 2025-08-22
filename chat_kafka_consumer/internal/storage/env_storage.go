package storage

import (
	"os"
)

type EnvStorage struct {
	DatabaseHost     string
	DatabasePort     string
	KafkaHost        string
	KafkaPort        string
	KafkaServicePort string
	KafkaTopic       string
	KafkaGroupId     string
}

var Env = &EnvStorage{}

func LoadEnv() {
	//err := godotenv.Load("../.env")
	//if err != nil {
	//	log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	//}

	Env.DatabaseHost = os.Getenv("SCYLLA_SERVICE_HOST")
	Env.DatabasePort = os.Getenv("SCYLLA_SERVICE_PORT")
	Env.KafkaHost = os.Getenv("KAFKA_HOST")
	Env.KafkaPort = os.Getenv("KAFKA_PORT")
	Env.KafkaServicePort = os.Getenv("KAFKA_SERVICE_PORT")
	Env.KafkaTopic = os.Getenv("KAFKA_TOPIC")
	Env.KafkaGroupId = os.Getenv("KAFKA_GROUP_ID")
}
