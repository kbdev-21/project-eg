package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"backend/src/shared"
	"strings"

	"github.com/gofiber/fiber/v3"
)

func InitHttpRoutes(a *fiber.App, q *db.Queries, conf config.Config) {
	a.Get("/api/me", authMiddleware(q, conf), func(ctx fiber.Ctx) error {
		u := ctx.Locals("currentUser").(db.User)

		return ctx.JSON(u)
	})

	a.Get("/api/users/:id", func(ctx fiber.Ctx) error {
		id, err := shared.ParseStringToUuid(ctx.Params("id", ""))
		if err != nil {
			return ctx.SendStatus(400)
		}
		user, err := q.GetUserById(ctx, id)
		if err != nil {
			return ctx.SendStatus(404)
		}
		return ctx.JSON(user)
	})
}

func authMiddleware(q *db.Queries, conf config.Config) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return ctx.SendStatus(401)
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := domain.VerifyToken(token, conf.JwkSet)
		if err != nil {
			return ctx.SendStatus(401)
		}

		currentUser, err := domain.SyncUserFromTokenPayload(ctx, q, payload)
		if err != nil {
			return ctx.SendStatus(401)
		}

		ctx.Locals("currentUser", currentUser)

		return ctx.Next()
	}
}
