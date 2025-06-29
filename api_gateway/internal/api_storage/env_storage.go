package api_storage

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"results/errs"
)

type EnvStorage struct {
	ApiPort   string
	AuthHost  string
	AuthPort  string
	DebugMode bool
}

var Env = &EnvStorage{}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	}

	Env.AuthHost = os.Getenv("AUTH_HOST")
	Env.ApiPort = os.Getenv("API_PORT")
	Env.AuthPort = os.Getenv("AUTH_PORT")
	Env.DebugMode = os.Getenv("DEBUG_MODE") == "true"
}
