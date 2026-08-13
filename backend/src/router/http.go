package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/shared"

	"github.com/gofiber/fiber/v3"
)

func InitHttpRoutes(a *fiber.App, q *db.Queries, conf config.Config) {
	a.Get("/api/me", authMiddleware(q, conf), getMeHandler())
	a.Get("/api/users/:id", getUserByIdHandler(q))
}

func getMeHandler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		u := ctx.Locals("currentUser").(db.User)
		return ctx.JSON(u)
	}
}

func getUserByIdHandler(q *db.Queries) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		id, err := shared.ParseStringToUuid(ctx.Params("id", ""))
		if err != nil {
			return ctx.SendStatus(400)
		}
		user, err := q.GetUserById(ctx, id)
		if err != nil {
			return ctx.SendStatus(404)
		}
		return ctx.JSON(user)
	}
}