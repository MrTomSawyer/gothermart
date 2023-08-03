package middleware

import (
	"time"

	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/gofiber/fiber/v2"
)

func (m *Middlewares) LogReqResInfo(c *fiber.Ctx) error {
	startTime := time.Now()
	c.Next()
	latency := time.Since(startTime)

	logger.Log.Infof("Method: %s URI: %s Latency: %v", c.Method(), c.Context().Path(), latency)
	return nil
}
