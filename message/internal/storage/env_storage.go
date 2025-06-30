package storage

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"results/errs"
)

type EnvStorage struct {
	MessagePort string
	JwtSecret   string
}

var Env = &EnvStorage{}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	}

	Env.MessagePort = os.Getenv("MESSAGE_PORT")
	Env.JwtSecret = os.Getenv("JWT_SECRET")
}
