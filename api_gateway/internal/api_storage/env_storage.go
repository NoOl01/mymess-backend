package api_storage

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type EnvStorage struct {
	ApiPort   string
	AuthPort  string
	DebugMode bool
}

var Env = &EnvStorage{}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Env.ApiPort = os.Getenv("API_PORT")
	Env.AuthPort = os.Getenv("AUTH_PORT")
	Env.DebugMode = os.Getenv("DEBUG_MODE") == "true"
}
