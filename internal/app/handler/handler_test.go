package handler

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/config"
	"github.com/MrTomSawyer/loyalty-system/internal/app/repository/mocks"
	"github.com/MrTomSawyer/loyalty-system/internal/app/service"
	"github.com/gofiber/fiber/v2"
)

func createTestServer(m *mocks.MockUserRepository) *fiber.App {
	app := fiber.New()

	cfg := config.NewConfig()
	cfg.UserTableName = "users"

	var auth *service.AuthService
	userService := service.NewUserService(cfg, m, auth)
	handler := NewHandler(userService)

	api := app.Group("/api")
	user := api.Group("/user")
	{
		user.Post("/register", handler.CreateUser)
	}

	return app
}
