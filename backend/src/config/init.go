package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var PostgresConnectionUrl string
var AppJwkSet JwkSet

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

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	PostgresConnectionUrl = os.Getenv("POSTGRES_CONNECTION_URL")
	AppJwkSet = fetchJwkSet(os.Getenv("JWKS_URL"))
}

func fetchJwkSet(url string) JwkSet {
	var jwkSet JwkSet

	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		log.Fatal("Failed to fetch JWKS")
	}

	if err := json.NewDecoder(res.Body).Decode(&jwkSet); err != nil {
		log.Fatal(err)
	}

	return jwkSet
}