package storage

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type EnvStorage struct {
	AuthPort     string
	DatabasePort string
	JwtSecret    string
}

var Env = &EnvStorage{}

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	Env.AuthPort = os.Getenv("AUTH_PORT")
	Env.DatabasePort = os.Getenv("DATABASE_PORT")
	Env.JwtSecret = os.Getenv("JWT_SECRET")
}
