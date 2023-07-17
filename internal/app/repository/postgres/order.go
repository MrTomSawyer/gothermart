package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/sqlerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	dbPool    *pgxpool.Pool
	ctx       context.Context
	tableName string
}

func NewOrderRepository(ctx context.Context, dbPool *pgxpool.Pool, tableName string) *OrderRepository {
	return &OrderRepository{
		dbPool:    dbPool,
		ctx:       ctx,
		tableName: tableName,
	}
}

func (o *OrderRepository) CreateOrder(order entity.Order) error {
	query := fmt.Sprintf("SELECT * FROM %s WHERE order_num=$1", o.tableName)
	row := o.dbPool.QueryRow(o.ctx, query, order.OrderID)

	var ord entity.Order
	err := row.Scan(&ord.ID, &ord.UserID, &ord.OrderID, &ord.Status, &ord.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log.Errorf("No rows for order_num %s found", order.OrderID)
		} else {
			return err
		}
	}

	if ord.OrderID != "" {
		if ord.UserID == order.UserID {
			logger.Log.Infof("This order has laready been uploaded by you")
			return sqlerr.ErrUploadedBySameUser
		}
		logger.Log.Infof("This order has already been uploaded by another user")
		return sqlerr.ErrUploadedByAnotherUser
	}

	query = fmt.Sprintf("INSERT INTO %s (user_id, order_num, order_status, created_at) VALUES ($1, $2, $3, $4) RETURNING id", o.tableName)
	row = o.dbPool.QueryRow(o.ctx, query, order.UserID, order.OrderID, order.Status, order.CreatedAt)

	var id int
	err = row.Scan(&id)
	if err != nil {
		logger.Log.Errorf("failed to create order %s", order.OrderID)
		return err
	}

	return nil
}

func (o *OrderRepository) GetAllOrders(userID int) ([]models.Order, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id=$1 ORDER BY to_timestamp(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SSOF') DESC", o.tableName)
	rows, err := o.dbPool.Query(o.ctx, query, userID)
	if err != nil {
		logger.Log.Errorf("failed to query rows for all orders: %v", err)
		return nil, err
	}

	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderID, &order.Accrual, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Errorf("failed to scan rows for all orders: %v", err)
		return nil, err
	}

	return orders, nil
}
