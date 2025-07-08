package storage

import (
	"os"
)

type EnvStorage struct {
	Email           string
	EmailPassword   string
	SmthHost        string
	SmtpPort        string
	SmtpServicePort string
}

var Env = &EnvStorage{}

func LoadEnv() {
	//err := godotenv.Load("../.env")
	//if err != nil {
	//	log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	//}

	Env.Email = os.Getenv("EMAIL")
	Env.EmailPassword = os.Getenv("EMAIL_PASSWORD")
	Env.SmthHost = os.Getenv("SMTP_HOST")
	Env.SmtpPort = os.Getenv("SMTP_PORT")
	Env.SmtpServicePort = os.Getenv("SMTP_SERVICE_PORT")
}
