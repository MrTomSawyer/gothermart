package service

import (
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/ordererr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/config"
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/interfaces"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"strconv"
	"strings"
	"time"
)

type OrderService struct {
	cfg             *config.Config
	orderRepository interfaces.OrderRepository
}

func NewOrderService(cfg *config.Config, orderRepository interfaces.OrderRepository) *OrderService {
	return &OrderService{
		cfg:             cfg,
		orderRepository: orderRepository,
	}
}

func (o *OrderService) CreateOrder(order models.Order) error {
	isValid := validateOrderID(order.OrderID)
	if !isValid {
		logger.Log.Errorf("%s has invalid format", order.OrderID)
		return ordererr.ErrWrongOrderID
	}
	orderEntity := entity.Order{
		UserID:    order.UserID,
		OrderID:   order.OrderID,
		Status:    "NEW",
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	err := o.orderRepository.CreateOrder(orderEntity)
	return err
}

func (o *OrderService) GetAllOrders(userID int) ([]models.Order, error) {
	orders, err := o.orderRepository.GetAllOrders(userID)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func validateOrderID(id string) bool {
	orderID := strings.ReplaceAll(id, " ", "")
	sum := 0
	shouldDouble := len(orderID)%2 == 0

	for i := 0; i <= len(orderID)-1; i++ {
		digit, err := strconv.Atoi(string(orderID[i]))
		if err != nil {
			return false
		}

		if shouldDouble {
			digit *= 2

			if digit > 9 {
				digit -= 9
			}
		}

		sum += digit
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}
