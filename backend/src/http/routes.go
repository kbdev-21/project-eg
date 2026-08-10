package http

import (
	"backend/src/db"
	"log"

	"github.com/gofiber/fiber/v3"
)

func InitRoutes(a *fiber.App, q *db.Queries) {
	initPublicRoutes(a, q)
	initAuthRoutes(a, q)
}

func initPublicRoutes(a *fiber.App, q *db.Queries) {
	pub := a.Group("")

	pub.Get("/api/hello", func(c fiber.Ctx) error {
		return c.JSON("Hello World")
	})
}

func initAuthRoutes(a *fiber.App, q *db.Queries) {
	auth := a.Group("")
	auth.Use(authMiddleware(q))

	auth.Get("/api/me", func(c fiber.Ctx) error {
		test, _ := c.Locals("test").(string)
		log.Println(test)
		return c.JSON("Hello World")
	})
}