package storage

import (
	"os"
)

type EnvStorage struct {
	DatabaseHost string
	DatabasePort string
	ProfileHost  string
	ProfilePort  string
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
	Env.ProfileHost = os.Getenv("PROFILE_SERVICE_HOST")
	Env.ProfilePort = os.Getenv("PROFILE_SERVICE_PORT")
	Env.JwtSecret = os.Getenv("JWT_SECRET")
}
