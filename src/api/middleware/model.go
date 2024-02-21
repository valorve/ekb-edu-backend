package middleware

import "github.com/gofiber/fiber/v2"

var TokenRequired fiber.Handler
var JwtSecret string
