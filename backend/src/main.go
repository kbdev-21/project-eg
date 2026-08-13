package main

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/router"
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	dbPool, err := pgxpool.New(context.Background(), config.PostgresConnectionUrl)
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

	router.InitPublicRoutes(app, queries)
	router.InitAuthRoutes(app, queries)
	router.InitWsRoute(app, queries)

	log.Fatal(app.Listen(":3000"))
}
