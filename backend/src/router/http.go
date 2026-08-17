package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"backend/src/shared"

	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitHttpRoutes(a *fiber.App, q *db.Queries, conf config.Config, pool *pgxpool.Pool) {
	a.Get("/api/me", authMiddleware(q, conf, pool), getMeHandler())
	a.Get("/api/users/:id", getUserByIdHandler(q, pool))
}

func getMeHandler() fiber.Handler {
	return func(ctx fiber.Ctx) error {
		u := ctx.Locals("currentUser").(domain.User)
		return ctx.JSON(u)
	}
}

func getUserByIdHandler(q *db.Queries, p *pgxpool.Pool) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		id, err := shared.ParseStringToUuid(ctx.Params("id", ""))
		if err != nil {
			return ctx.SendStatus(400)
		}
		user, err := domain.GetUserById(domain.CreateDbExecDeps(ctx, q, p), id)
		if err != nil {
			return ctx.SendStatus(404)
		}
		return ctx.JSON(user)
	}
}