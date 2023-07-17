package handler

import "github.com/MrTomSawyer/loyalty-system/internal/app/interfaces"

type Handler struct {
	userService       interfaces.UserService
	orderService      interfaces.OrderService
	withdrawalService interfaces.WithdrawalService
}

func NewHandler(userService interfaces.UserService, orderService interfaces.OrderService, withdrawalService interfaces.WithdrawalService) *Handler {
	return &Handler{
		userService:       userService,
		orderService:      orderService,
		withdrawalService: withdrawalService,
	}
}
