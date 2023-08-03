package interfaces

import "github.com/gofiber/fiber/v2"

type Handler interface {
	Login(c *fiber.Ctx) error
	CreateUser(c *fiber.Ctx) error
	CreateOrder(c *fiber.Ctx) error
	GetAllOrders(c *fiber.Ctx) error
	GetUserBalance(c *fiber.Ctx) error
	Withdraw(c *fiber.Ctx) error
	GetWithdrawals(c *fiber.Ctx) error
}
