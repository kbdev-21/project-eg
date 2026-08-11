package http

import (
	"backend/src/db"
	"github.com/gofiber/fiber/v3"
)

func InitPublicRoutes(a *fiber.App, q *db.Queries) {
	pub := a.Group("")

	pub.Get("/api/hello", func(c fiber.Ctx) error {
		return c.JSON("Hello World")
	})
}

func InitAuthRoutes(a *fiber.App, q *db.Queries) {
	auth := a.Group("")
	auth.Use(authMiddleware(q))

	auth.Get("/api/me", func(c fiber.Ctx) error {
		u := c.Locals("currentUser").(db.User)

		return c.JSON(u)
	})
}