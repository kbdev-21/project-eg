package router

import (
	"backend/src/config"
	"backend/src/domain"
	"backend/src/shared"

	"github.com/gofiber/fiber/v3"
)

func InitHttpRoutes(fib *fiber.App, a *domain.AppState, cfg config.Config) {
	fib.Get("/api/me", authMiddleware(a, cfg), getMeHandler())
	fib.Get("/api/users/:id", getUserByIdHandler(a))
}

func getMeHandler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		u := ctx.Locals("currentUser").(domain.User)
		return ctx.JSON(u)
	}
}

func getUserByIdHandler(a *domain.AppState) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		id, err := shared.ParseStringToUuid(ctx.Params("id", ""))
		if err != nil {
			return ctx.SendStatus(400)
		}
		user, err := a.GetUserById(ctx, id)
		if err != nil {
			return ctx.SendStatus(404)
		}
		return ctx.JSON(user)
	}
}
