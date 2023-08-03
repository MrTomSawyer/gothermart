package handler

import (
	"errors"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/ordererr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/sqlerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

func (h *Handler) CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	order := models.Order{
		UserID:  strconv.Itoa(userID),
		OrderID: string(c.Body()),
	}

	err := h.orderService.CreateOrder(order)
	if err != nil {
		if errors.Is(err, ordererr.ErrWrongOrderID) {
			return c.Status(http.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		if errors.Is(err, sqlerr.ErrUploadedBySameUser) {
			return c.Status(http.StatusOK).JSON(fiber.Map{
				"message": err.Error(),
			})
		}
		if errors.Is(err, sqlerr.ErrUploadedByAnotherUser) {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"message": err.Error(),
			})
		}

		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(http.StatusAccepted).JSON(fiber.Map{
		"message": "New order number has been received and is being processed",
	})
}

func (h *Handler) GetAllOrders(c *fiber.Ctx) error {
	userID := c.Locals("userID").(int)

	orders, err := h.orderService.GetAllOrders(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": err.Error(),
		})
	}
	if len(orders) == 0 {
		return c.SendStatus(http.StatusNoContent)
	}

	return c.Status(http.StatusOK).JSON(orders)
}
