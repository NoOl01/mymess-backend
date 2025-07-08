package storage

import (
	"os"
)

type EnvStorage struct {
	AuthPort     string
	DatabaseHost string
	DatabasePort string
	SmtpHost     string
	SmtpPort     string
	JwtSecret    string
}

var Env = &EnvStorage{}

func LoadEnv() {
	//err := godotenv.Load("../.env")
	//if err != nil {
	//	log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	//}

	Env.DatabaseHost = os.Getenv("DB_SERVICE_HOST")
	Env.DatabasePort = os.Getenv("DB_SERVICE_PORT")
	Env.AuthPort = os.Getenv("AUTH_PORT")
	Env.SmtpHost = os.Getenv("SMTP_SERVICE_HOST")
	Env.SmtpPort = os.Getenv("SMTP_SERVICE_PORT")
	Env.JwtSecret = os.Getenv("JWT_SECRET")
}
