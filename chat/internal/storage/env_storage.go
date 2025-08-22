package storage

import (
	"os"
)

type EnvStorage struct {
	KafkaHost   string
	KafkaPort   string
	MessagePort string
	JwtSecret   string
}

var Env = &EnvStorage{}

func LoadEnv() {
	//if err := godotenv.Load("../.env"); err != nil {
	//	log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	//}

	Env.KafkaHost = os.Getenv("KAFKA_HOST")
	Env.KafkaPort = os.Getenv("KAFKA_PORT")
	Env.MessagePort = os.Getenv("CHAT_PORT")
	Env.JwtSecret = os.Getenv("JWT_SECRET")
}
