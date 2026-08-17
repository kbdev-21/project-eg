package router

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"strings"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"github.com/jackc/pgx/v5/pgxpool"
)

func authMiddleware(q *db.Queries, conf config.Config, p *pgxpool.Pool) fiber.Handler {
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

		currentUser, err := domain.SyncUserFromTokenPayload(domain.CreateDbExec(ctx, q, p), payload)
		if err != nil {
			return ctx.SendStatus(401)
		}

		ctx.Locals("currentUser", currentUser)

		return ctx.Next()
	}
}

func wsAuthMiddleware(q *db.Queries, conf config.Config, p *pgxpool.Pool) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			token := ctx.Query("token", "")

			payload, err := domain.VerifyToken(token, conf.JwkSet)
			if err != nil {
				return ctx.SendStatus(401)
			}

			currentUser, err := domain.SyncUserFromTokenPayload(domain.CreateDbExec(ctx, q, p), payload)
			if err != nil {
				return ctx.SendStatus(401)
			}

			ctx.Locals("currentUser", currentUser)

			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
