package storage

import (
	"os"
)

type EnvStorage struct {
	ScyllaHost        string
	ScyllaServicePort string
	ScyllaSpaceKey    string
}

var Env = &EnvStorage{}

func LoadEnv() {
	//if err := godotenv.Load("../.env"); err != nil {
	//	log.Fatalf("%s: %s", errs.FailedLoadEnvFile, err.Error())
	//}

	Env.ScyllaHost = os.Getenv("SCYLLA_HOST")
	Env.ScyllaServicePort = os.Getenv("SCYLLA_SERVICE_PORT")
	Env.ScyllaSpaceKey = os.Getenv("SCYLLA_SPACE_KEY")
}
