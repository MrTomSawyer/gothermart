package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (m *Middlewares) AuthMiddleware(c *fiber.Ctx) error {
	token := c.Cookies("Authorization")

	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing token cookie",
		})
	}

	userID, err := m.auth.GetUserID(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	c.Locals("userID", userID)
	c.Next()
	return nil
}
