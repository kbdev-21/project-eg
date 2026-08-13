package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"backend/src/shared"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func InitPublicRoutes(a *fiber.App, q *db.Queries) {
	a.Get("/api/users/:id", func(c fiber.Ctx) error {
		id, err := shared.ParseStringToUuid(c.Params("id", "")) 
		if err != nil {
			return c.SendStatus(400)
		}
		user, err := q.GetUserById(c, id)
		if err != nil {
			return c.SendStatus(404)
		}
		return c.JSON(user)
	})
}

func InitAuthRoutes(a *fiber.App, q *db.Queries) {
	a.Get("/api/me", authMiddleware(q), func(c fiber.Ctx) error {
		u := c.Locals("currentUser").(db.User)

		return c.JSON(u)
	})
}

func authMiddleware(q *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.SendStatus(401)
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := domain.VerifyToken(token, config.AppJwkSet)
		if err != nil {
			return c.SendStatus(401)
		}

		currentUser, err := domain.SyncUserFromTokenPayload(c, q, payload)
		if err != nil {
			return c.SendStatus(401)
		}

		c.Locals("currentUser", currentUser)

		return c.Next()
	}
}