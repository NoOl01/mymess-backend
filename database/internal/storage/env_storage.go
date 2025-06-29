package storage

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"results/errs"
)

type EnvStorage struct {
	DbPort string
	DbHost string
	DbUser string
	DbPass string
	DbName string
	DbSsl  string
}

var Env = &EnvStorage{}

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	}

	Env.DbPort = os.Getenv("DB_PORT")
	Env.DbHost = os.Getenv("DB_HOST")
	Env.DbUser = os.Getenv("DB_USER")
	Env.DbPass = os.Getenv("DB_PASS")
	Env.DbName = os.Getenv("DB_NAME")
	Env.DbSsl = os.Getenv("DB_SSL")
}
