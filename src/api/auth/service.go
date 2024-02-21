package auth

import (
	"crypto/sha256"
	"ekb-edu/src/api/middleware"
	"ekb-edu/src/database/storage"
	"encoding/hex"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

func getPasswordHash(password string) string {
	sum := sha256.Sum256([]byte(password))
	return hex.EncodeToString(sum[:])
}

func getClaims(user storage.EeUser) jwt.MapClaims {
	return jwt.MapClaims{
		"name":  user.Username,
		"admin": user.Username == "izke",
		"id":    user.UserID,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}
}

func register(c *fiber.Ctx) error {
	userInfo := User{}

	if err := c.BodyParser(&userInfo); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "cannot parse auth data"})
	}

	// Проверяем, существует ли уже пользователь с таким же username или email
	var count int64
	storage.DB.Model(&storage.EeUser{}).Where("username = ? OR email = ?", userInfo.Username, userInfo.Email).Count(&count)
	if count > 0 {
		return gorm.ErrRecordNotFound
	}

	// Создаём и сохраняем пользователя
	user := storage.EeUser{
		Username:     userInfo.Username,
		Email:        userInfo.Email,
		PasswordHash: getPasswordHash(userInfo.Password),
	}

	result := storage.DB.Create(&user)
	if result.Error != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "failed to register"})
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, getClaims(user))

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(middleware.JwtSecret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create token"})
	}

	return c.JSON(fiber.Map{"token": t})
}

func login(c *fiber.Ctx) error {
	userInfo := User{}

	if err := c.BodyParser(&userInfo); err != nil {
		return c.JSON(fiber.Map{"error": err})
	}

	// Проверяем, существует ли уже пользователь с таким же username или email
	var user storage.EeUser
	result := storage.DB.Model(&storage.EeUser{}).Where("username = ? AND password_hash = ?", userInfo.Username, getPasswordHash(userInfo.Password)).First(&user)

	if result.Error != nil {
		return result.Error
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, getClaims(user))

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(middleware.JwtSecret))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func restricted(c *fiber.Ctx) error {
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	name := claims["name"].(string)
	return c.SendString("Welcome " + name)
}

func RegisterService(app fiber.Router) {
	g := app.Group("/auth")
	g.Post("/register", register)
	g.Post("/login", login)

	g.Get("/restricted", middleware.TokenRequired, restricted)
}
