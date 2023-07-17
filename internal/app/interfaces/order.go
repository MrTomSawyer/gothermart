package interfaces

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
)

type OrderRepository interface {
	CreateOrder(order entity.Order) error
	GetAllOrders(userID int) ([]models.Order, error)
}

type OrderService interface {
	CreateOrder(order models.Order) error
	GetAllOrders(userID int) ([]models.Order, error)
}
