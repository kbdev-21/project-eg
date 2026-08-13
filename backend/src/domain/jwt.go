package domain

import (
	"backend/src/config"
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/golang-jwt/jwt/v5"
)

type TokenPayload struct {
	jwt.RegisteredClaims

	Sub          string `json:"sub"`
	IsAnonymous  bool   `json:"is_anonymous"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	UserMetadata struct {
		AvatarUrl string `json:"avatar_url"`
		Email     string `json:"email"`
		FullName  string `json:"full_name"`
		Name      string `json:"name"`
		Picture   string `json:"picture"`
	} `json:"user_metadata"`
}

func VerifyToken(token string, jwkSet config.JwkSet) (TokenPayload, error) {
	var payload TokenPayload

	_, err := jwt.ParseWithClaims(token, &payload, getPublicKey(jwkSet))

	return payload, err
}

func getPublicKey(jwkSet config.JwkSet) jwt.Keyfunc {
	return func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodES256 {
			return nil, fmt.Errorf("unexpected signing method")
		}

		kid := t.Header["kid"].(string)

		for _, key := range jwkSet.Keys {
			if key.Kid == kid {
				x, _ := base64.RawURLEncoding.DecodeString(key.X)
				y, _ := base64.RawURLEncoding.DecodeString(key.Y)

				return &ecdsa.PublicKey{
					Curve: elliptic.P256(),
					X:     new(big.Int).SetBytes(x),
					Y:     new(big.Int).SetBytes(y),
				}, nil
			}
		}

		return nil, fmt.Errorf("key not found")
	}
}

func printTokenPayload(token string) {
	t, _, err := jwt.NewParser().ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		fmt.Println(err)
		return
	}

	data, _ := json.MarshalIndent(t.Claims, "", "  ")
	fmt.Println(string(data))
}
