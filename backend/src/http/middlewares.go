package http

import (
	"backend/src/db"

	"github.com/gofiber/fiber/v3"
)

func authMiddleware(q *db.Queries) fiber.Handler {
	return func(c fiber.Ctx) error {
		c.Locals("test", "test")
		return c.Next()
	}
}