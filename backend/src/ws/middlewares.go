package ws

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain/auth"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func wsMiddleware(q *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			token := c.Query("token", "")

			payload, err := auth.VerifyToken(token, config.AppJwkSet)
			if err != nil {
				return c.SendStatus(401)
			}

			currentUser, err := auth.SyncUserFromTokenPayload(c, q, payload)
			if err != nil {
				return c.SendStatus(401)
			}

			c.Locals("currentUser", currentUser)

			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
