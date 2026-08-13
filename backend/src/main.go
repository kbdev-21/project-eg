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
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := pgxpool.New(context.Background(), conf.PostgresConnectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	queries := db.New(dbPool)

	app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(ctx fiber.Ctx) error {
		return ctx.JSON("Welcome to Project EG")
	})

	router.InitHttpRoutes(app, queries, conf)
	router.InitWsRoute(app, queries, conf)

	log.Fatal(app.Listen(":3000"))
}
