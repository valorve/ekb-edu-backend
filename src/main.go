package main

import (
	"ekb-edu/src/api/auth"
	"ekb-edu/src/api/courses"
	"ekb-edu/src/api/courses/lessons"
	"ekb-edu/src/api/middleware"
	"ekb-edu/src/database/config"
	"ekb-edu/src/database/storage"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	storage.Connect(&cfg.Postgres)
	middleware.InitializeJWT(&cfg.Jwt)

	app := fiber.New()
	{
		config := cors.ConfigDefault
		config.AllowCredentials = true
		app.Use(cors.New(config))
	}

	app.Use(logger.New(logger.Config{
		Format:     "[${time}] ${status} - ${latency} ${method} ${reqHeader:X-Forwarded-For} ${reqHeader:X-User-Id} ${reqHeader:X-User-Login} ${path} ${queryParams}\n",
		TimeFormat: "2006-01-02 15:04:05",
	}))

	v1 := app.Group("/v1")

	auth.RegisterService(v1)
	courses.RegisterService(v1)
	lessons.RegisterService(v1)
	app.Listen(fmt.Sprintf(":%d", cfg.Web.Port))
}
