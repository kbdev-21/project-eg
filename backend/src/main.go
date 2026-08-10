package main

import (
	"backend/src/config"
	"backend/src/db"
	"backend/src/http"
	"context"
	"log"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/jackc/pgx/v5"
)

func main() {
	env := config.LoadEnvVariables()

	dbConn, err := pgx.Connect(context.Background(), env.PostgresConnectionUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer dbConn.Close(context.Background())

	queries := db.New(dbConn)

    app := fiber.New()

	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c fiber.Ctx) error {
		return c.JSON("Welcome to Project EG")
	})

    http.InitRoutes(app, queries)

    log.Fatal(app.Listen(":3000"))
}