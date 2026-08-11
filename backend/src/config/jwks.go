package config

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type JwkSet struct {
	Keys []Jwk `json:"keys"`
}

type Jwk struct {
	X      string   `json:"x"`
	Y      string   `json:"y"`
	Alg    string   `json:"alg"`
	Crv    string   `json:"crv"`
	Ext    bool     `json:"ext"`
	Kid    string   `json:"kid"`
	Kty    string   `json:"kty"`
	KeyOps []string `json:"key_ops"`
}

var AppJwkSet JwkSet

func FetchJwkSet() {
	jwksUrl := os.Getenv("JWKS_URL")

	var jwkSet JwkSet

	res, err := http.Get(jwksUrl)
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

	AppJwkSet = jwkSet
}
