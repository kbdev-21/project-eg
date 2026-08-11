package main

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/http"
	"backend/src/ws"
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config.LoadEnv()
	config.FetchJwkSet()

	dbPool, err := pgxpool.New(context.Background(), config.Env.PostgresConnectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	queries := db.New(dbPool)

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON("Welcome to Project EG")
	})

	http.InitPublicRoutes(app, queries)
	http.InitAuthRoutes(app, queries)
	ws.InitWsRoutes(app, queries)
	

	log.Fatal(app.Listen(":3000"))
}
