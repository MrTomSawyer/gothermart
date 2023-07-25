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

//serv
//mockgen -source=C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/interfaces/order.go -destination=C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/repository/mocks/mock_orderserv.go -package=mocks C:/Users/lette/OneDrive/Документы/Projects/ygo/go-musthave-diploma-tpl/internal/app/service/order.go OrderService
