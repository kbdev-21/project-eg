package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvVariables struct {
	PostgresConnectionUrl string
}

func LoadEnvVariables() EnvVariables {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env")
	}

	return EnvVariables{
		PostgresConnectionUrl: os.Getenv("POSTGRES_CONNECTION_URL"),
	}
}
