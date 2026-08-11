package http

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain/auth"
	"log"
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
			log.Println(err)
			return c.SendStatus(401)
		}

		currentUser, err := auth.SyncUserFromTokenPayload(c, q, payload)
		if err != nil {
			log.Println(err)
			return c.SendStatus(401)
		}

		c.Locals("currentUser", currentUser)

		return c.Next()
	}
}
