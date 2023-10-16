package middleware

import (
	"fmt"
	"strings"

	"github.com/akhil-is-watching/authservice/config"
	"github.com/akhil-is-watching/authservice/repository"
	"github.com/akhil-is-watching/authservice/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func DeserializeUser(c *fiber.Ctx) error {
	var tokenString string
	authorization := c.Get("Authorization")
	config, _ := config.LoadConfig(".")

	if strings.HasPrefix(authorization, "Bearer ") {
		tokenString = strings.TrimPrefix(authorization, "Bearer ")
	} else if c.Cookies("token") != "" {
		tokenString = c.Cookies("token")
	}

	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "You are not logged in"})
	}

	tokenByte, err := jwt.Parse(tokenString, func(jwtToken *jwt.Token) (interface{}, error) {
		if _, ok := jwtToken.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %s", jwtToken.Header["alg"])
		}

		return []byte(config.JwtSecret), nil
	})

	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": fmt.Sprintf("invalidate token: %v", err)})
	}

	claims, ok := tokenByte.Claims.(jwt.MapClaims)
	if !ok || !tokenByte.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"status": "fail", "message": "invalid token claim"})

	}

	userRepo := repository.NewUserRepository(storage.GetDB())
	user, err := userRepo.Get(fmt.Sprint(claims["sub"]))
	if err != nil {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status": "fail", "message": "the user belonging to this token no logger exists"})
	}

	c.Locals("user", user)

	return c.Next()
}
