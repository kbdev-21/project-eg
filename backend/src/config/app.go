package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PostgresConnectionUrl string
	JwkSet JwkSet
}

type JwkSet struct {
	Keys []struct {
		X      string   `json:"x"`
		Y      string   `json:"y"`
		Alg    string   `json:"alg"`
		Crv    string   `json:"crv"`
		Ext    bool     `json:"ext"`
		Kid    string   `json:"kid"`
		Kty    string   `json:"kty"`
		KeyOps []string `json:"key_ops"`
	} `json:"keys"`
}

func LoadConfig() (Config, error) {
	err := godotenv.Load()
	if err != nil {
		return Config{}, err
	}

	pgConnUrl := os.Getenv("POSTGRES_CONNECTION_URL")
	jwkSet, err := fetchJwkSet(os.Getenv("JWKS_URL"))
	if err != nil {
		return Config{}, err
	}

	return Config{
		PostgresConnectionUrl: pgConnUrl,
		JwkSet: jwkSet,
	}, nil
}

func fetchJwkSet(url string) (JwkSet, error) {
	var jwkSet JwkSet

	res, err := http.Get(url)
	if err != nil {
		return jwkSet, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return jwkSet, fmt.Errorf("failed fetching JWKS")
	}

	if err := json.NewDecoder(res.Body).Decode(&jwkSet); err != nil {
		return jwkSet, err
	}

	return jwkSet, nil
}