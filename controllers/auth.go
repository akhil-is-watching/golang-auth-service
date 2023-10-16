package controllers

import (
	"time"

	"github.com/akhil-is-watching/authservice/config"
	"github.com/akhil-is-watching/authservice/models"
	"github.com/akhil-is-watching/authservice/repository"
	"github.com/akhil-is-watching/authservice/storage"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	params := &models.UserCreate{}
	if err := c.BodyParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	if err := params.Validate(); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	user, err := params.Convert()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	userRepo := repository.NewUserRepository(storage.GetDB())
	if err := userRepo.Create(user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func SignIn(c *fiber.Ctx) error {
	params := &models.UserSignIn{}

	if err := c.BodyParser(params); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	userRepo := repository.NewUserRepository(storage.GetDB())
	user, err := userRepo.Get(params.Email)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": err.Error(),
		})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(params.Password))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"msg": "Invalid Email or Password",
		})
	}

	tokenByte := jwt.New(jwt.SigningMethodHS256)
	now := time.Now().UTC()
	claims := tokenByte.Claims.(jwt.MapClaims)
	config, _ := config.LoadConfig(".")
	claims["sub"] = user.ID
	claims["exp"] = now.Add(config.JwtExpiresIn).Unix()
	claims["iat"] = now.Unix()
	claims["nbf"] = now.Unix()

	tokenString, err := tokenByte.SignedString([]byte(config.JwtSecret))

	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"msg": "Generating JWT Token failed",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "token",
		Value:    tokenString,
		Path:     "/",
		MaxAge:   config.JwtMaxAge * 60,
		Secure:   false,
		HTTPOnly: true,
		Domain:   "localhost",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": tokenString})
}

func LogoutUser(c *fiber.Ctx) error {
	expired := time.Now().Add(-time.Hour * 24)
	c.Cookie(&fiber.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expired,
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success"})
}
