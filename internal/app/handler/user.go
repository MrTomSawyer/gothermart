package handler

import (
	"encoding/json"
	"errors"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/autherr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/sqlerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *Handler) Login(c *fiber.Ctx) error {
	var user models.User
	err := json.Unmarshal(c.Body(), &user)

	if user.Login == "" || user.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Login or password must not be empty",
		})
	}

	if err != nil {
		logger.Log.Errorf("failed to unmarshal req body: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	jwt, err := h.userService.Login(user)
	if err != nil {
		if errors.Is(err, autherr.ErrWrongCredentials) || errors.Is(err, sqlerr.ErrNoRows) {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "Authorization",
		Value:    jwt,
		HTTPOnly: true,
	})
	return c.Status(http.StatusOK).JSON(nil)
}

func (h *Handler) CreateUser(c *fiber.Ctx) error {
	var user models.User
	err := json.Unmarshal(c.Body(), &user)

	if user.Login == "" || user.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Login or password must not be empty",
		})
	}

	if err != nil {
		logger.Log.Errorf("failed to unmarshal req body: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.userService.CreateUser(user)
	if err != nil {
		if errors.Is(err, sqlerr.ErrLoginConflict) {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		logger.Log.Errorf("failed to create a user: %v", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return nil
}

func (h *Handler) GetUserBalance(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	balance, err := h.userService.GetUserBalance(userID)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(balance)
}
