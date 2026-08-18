package main

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/domain"
	"backend/src/router"
	"backend/src/schedule"
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	dbPool, err := pgxpool.New(context.Background(), cfg.PostgresConnectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbPool.Close()

	queries := db.New(dbPool)

	fib := fiber.New()

	fib.Use(cors.New())
	fib.Use(logger.New())

	fib.Get("/", func(ctx fiber.Ctx) error {
		return ctx.JSON("Welcome to Project EG")
	})

	appState := domain.InitAppState(queries, dbPool)

	router.InitHttpRoutes(fib, appState, cfg)
	router.InitWsRoute(fib, appState, cfg)

	go schedule.EverySecond(appState)

	log.Fatal(fib.Listen(":3000"))
}
