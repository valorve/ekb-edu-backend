package middleware

import (
	"ekb-edu/src/database/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func hasClaim(c *fiber.Ctx, claim string) bool {
	user := c.Locals("user").(*jwt.Token)
	if user == nil {
		return false
	}

	claims := user.Claims.(jwt.MapClaims)
	return claims[claim].(bool)
}

func IsAdmin(c *fiber.Ctx) bool {
	return hasClaim(c, "admin")
}

func AdminRequired(c *fiber.Ctx) error {
	if !hasClaim(c, "admin") {
		return c.SendStatus(fiber.StatusForbidden)
	}

	return c.Next()
}

func InitializeJWT(cfg *config.Jwt) {
	JwtSecret = cfg.Secret

	TokenRequired = jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(JwtSecret)},
	})
}
