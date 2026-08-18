package router

import (
	"backend/src/config"
	"backend/src/domain"
	"strings"

	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
)

func authMiddleware(a *domain.AppState, cfg config.Config) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return ctx.SendStatus(401)
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")

		payload, err := domain.VerifyToken(token, cfg.JwkSet)
		if err != nil {
			return ctx.SendStatus(401)
		}

		currentUser, err := a.SyncUserFromTokenPayload(ctx, payload)
		if err != nil {
			return ctx.SendStatus(401)
		}

		ctx.Locals("currentUser", currentUser)

		return ctx.Next()
	}
}

func wsAuthMiddleware(a *domain.AppState, cfg config.Config) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(ctx) {
			token := ctx.Query("token", "")

			payload, err := domain.VerifyToken(token, cfg.JwkSet)
			if err != nil {
				return ctx.SendStatus(401)
			}

			currentUser, err := a.SyncUserFromTokenPayload(ctx, payload)
			if err != nil {
				return ctx.SendStatus(401)
			}

			ctx.Locals("currentUser", currentUser)

			return ctx.Next()
		}
		return fiber.ErrUpgradeRequired
	}
}
