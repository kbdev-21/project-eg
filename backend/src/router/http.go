package router

import (
	"backend/src/config"
	"backend/src/domain"
	"backend/src/shared"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func InitHttpRoutes(fib *fiber.App, a *domain.AppState, cfg config.Config) {
	fib.Get("/api/me", authMiddleware(a, cfg), getMeHandler())
	fib.Put("/api/me", authMiddleware(a, cfg), updateMeHandler(a))
	fib.Get("/api/users/:id", getUserByIdHandler(a))
}

func getMeHandler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		u := ctx.Locals("currentUser").(domain.User)
		return ctx.JSON(u)
	}
}

func updateMeHandler(a *domain.AppState) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var req domain.UpdateUserReq

		err := ctx.Bind().Body(&req)
		if err != nil {
			return ctx.SendStatus(400)
		}

		err = validator.New().Struct(req)
		if err != nil {
			return ctx.SendStatus(400)
		}

		u := ctx.Locals("currentUser").(domain.User)

		updatedUser, err := a.UpdateUserInfo(ctx, u.ID, req)
		if err != nil {
			return ctx.SendStatus(409)
		}

		return ctx.JSON(updatedUser)
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
