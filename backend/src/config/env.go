package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	PostgresConnectionUrl string
	JwksUrl               string
}

var Env EnvVariables

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	Env = EnvVariables{
		PostgresConnectionUrl: os.Getenv("POSTGRES_CONNECTION_URL"),
		JwksUrl:               os.Getenv("JWKS_URL"),
	}
}
