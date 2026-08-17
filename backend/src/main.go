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

	fib := fiber.New()

	fib.Use(cors.New())
	fib.Use(logger.New())

	fib.Get("/", func(ctx fiber.Ctx) error {
		return ctx.JSON("Welcome to Project EG")
	})

	appState := domain.InitAppState()

	router.InitHttpRoutes(fib, queries, conf, dbPool)
	router.InitWsRoute(fib, queries, conf, appState, dbPool)

	go schedule.EverySecond(appState, queries, dbPool)

	log.Fatal(fib.Listen(":3000"))
}
