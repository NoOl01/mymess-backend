package storage

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"results/errs"
)

type EnvStorage struct {
	AuthPort     string
	DatabaseHost string
	DatabasePort string
	JwtSecret    string
}

var Env = &EnvStorage{}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	}

	Env.DatabaseHost = os.Getenv("DATABASE_HOST")
	Env.AuthPort = os.Getenv("AUTH_PORT")
	Env.DatabasePort = os.Getenv("DATABASE_PORT")
	Env.JwtSecret = os.Getenv("JWT_SECRET")
}
