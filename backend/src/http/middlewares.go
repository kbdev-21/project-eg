package http

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain/auth"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func authMiddleware(q *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.SendStatus(401)
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := auth.VerifyToken(token, config.AppJwkSet)
		if err != nil {
			fmt.Println(err)
			return c.SendStatus(401)
		}

		fmt.Println(payload.Email)

		return c.Next()
	}
}
