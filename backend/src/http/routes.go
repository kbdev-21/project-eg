package http

import (
	"backend/src/db"
	"backend/src/shared"

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