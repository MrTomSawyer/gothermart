package handler

import (
	"encoding/json"
	"errors"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/ordererr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/withdrawerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func (h *Handler) Withdraw(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	var withdrawal models.Withdraw
	err := json.Unmarshal(c.Body(), &withdrawal)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	err = h.withdrawalService.Withdraw(withdrawal, userID)
	if err != nil {
		if errors.Is(err, withdrawerr.ErrLowBalance) {
			return c.Status(http.StatusPaymentRequired).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		if errors.Is(err, ordererr.ErrWrongOrderID) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "points has successfully been withdrawn",
	})
}

func (h *Handler) GetWithdrawals(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	withdrawals, err := h.withdrawalService.GetWithdrawals(userID)
	if err != nil {
		if errors.Is(err, withdrawerr.ErrNoWithdrawals) {
			return c.SendStatus(http.StatusNoContent)
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusOK).JSON(withdrawals)
}
