package interfaces

import "github.com/gofiber/fiber/v2"

type Middlewares interface {
	AuthMiddleware(c *fiber.Ctx) error
	LogReqResInfo(c *fiber.Ctx) error
}
