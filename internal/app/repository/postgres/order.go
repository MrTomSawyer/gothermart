package postgres

import (
	"context"
	"errors"
	"github.com/MrTomSawyer/loyalty-system/internal/app/apperrors/sqlerr"
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	dbPool *pgxpool.Pool
	ctx    context.Context
}

func NewOrderRepository(ctx context.Context, dbPool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{
		dbPool: dbPool,
		ctx:    ctx,
	}
}

func (o *OrderRepository) CreateOrder(order entity.Order) error {
	row := o.dbPool.QueryRow(o.ctx, "SELECT * FROM orders WHERE order_num=$1", order.OrderID)

	var ord entity.Order
	err := row.Scan(&ord.ID, &ord.UserID, &ord.OrderID, &ord.Accrual, &ord.Status, &ord.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			logger.Log.Infof("No rows for order_num %s found. Saving order...", order.OrderID)
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

	row = o.dbPool.QueryRow(o.ctx,
		"INSERT INTO orders (user_id, order_num, order_status, created_at) VALUES ($1, $2, $3, $4) RETURNING id",
		order.UserID, order.OrderID, order.Status, order.CreatedAt)

	var id int
	err = row.Scan(&id)
	if err != nil {
		logger.Log.Errorf("failed to create order %s", order.OrderID)
		return err
	}

	return nil
}

func (o *OrderRepository) UpdateOrderAccrual(order models.Order, orderID string, userID int) error {
	tx, err := o.dbPool.BeginTx(o.ctx, pgx.TxOptions{IsoLevel: pgx.Serializable})
	if err != nil {
		logger.Log.Infof("failed to begin transaction: %v", err)
		return err
	}

	_, err = tx.Exec(o.ctx, "UPDATE orders SET order_status=$1, accrual=$2 WHERE order_num=$3", order.Status, order.Accrual, orderID)
	if err != nil {
		if err := tx.Rollback(o.ctx); err != nil {
			logger.Log.Errorf("failed to rollback transaction: %v", err)
			return err
		}
		logger.Log.Errorf("failed to update order=%s: %v", orderID, err)
		return err
	}

	row := tx.QueryRow(o.ctx, "SELECT balance FROM users WHERE id=$1", userID)
	var user entity.User
	err = row.Scan(&user.Balance)
	if err != nil {
		if rollbackErr := tx.Rollback(o.ctx); rollbackErr != nil {
			logger.Log.Errorf("failed to rollback transaction: %v", rollbackErr)
			return err
		}
		logger.Log.Errorf("failed to find user with id=%d : %v", userID, err)
		return err
	}

	_, err = tx.Exec(o.ctx, "UPDATE users SET balance=balance+$1 WHERE id=$2", order.Accrual, userID)
	if err != nil {
		if rollbackErr := tx.Rollback(o.ctx); rollbackErr != nil {
			logger.Log.Errorf("failed to rollback transaction: %v", err)
			return err
		}
		logger.Log.Errorf("failed to upfate user with id=%d : %v", userID, err)
		return err
	}

	if err := tx.Commit(o.ctx); err != nil {
		logger.Log.Errorf("failed to commit transaction: %v", err)
	}
	logger.Log.Infof("Transaction successfully commited")

	return nil
}

func (o *OrderRepository) GetUnhandledOrders() ([]string, error) {
	var orderIDs []string
	rows, err := o.dbPool.Query(o.ctx, "SELECT order_num FROM orders WHERE order_status='NEW'")
	if err != nil {
		logger.Log.Errorf("failed to query rows for all orders with NEW status: %v", err)
		return []string{}, err
	}
	for rows.Next() {
		var ID string
		err := rows.Scan(&ID)
		if err != nil {
			logger.Log.Errorf("failder to stan order ID: %v", err)
			return []string{}, err
		}
		orderIDs = append(orderIDs, ID)
	}

	return orderIDs, nil
}

func (o *OrderRepository) GetOrderAndUserIDs(orderID string) (int, error) {
	row := o.dbPool.QueryRow(o.ctx, "SELECT user_id, order_status from orders WHERE order_num=$1", orderID)
	var userID int
	var orderStatus string
	err := row.Scan(&userID, &orderStatus)
	if err != nil {
		logger.Log.Errorf("no order with id=%s found: %v", orderID, err)
		return 0, err
	}
	if orderStatus == "INVALID" || orderStatus == "PROCESSED" {
		logger.Log.Errorf("order with id=%s has status of %s and can't be processed", orderID, orderStatus)
		return 0, err
	}

	return userID, nil
}

func (o *OrderRepository) GetAllOrders(userID int) ([]models.Order, error) {
	rows, err := o.dbPool.Query(o.ctx, "SELECT * FROM orders WHERE user_id=$1 ORDER BY to_timestamp(created_at, 'YYYY-MM-DD\"T\"HH24:MI:SS') DESC", userID)
	if err != nil {
		logger.Log.Errorf("failed to query rows for all orders: %v", err)
		return nil, err
	}

	var orders []models.Order
	for rows.Next() {
		order := models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.OrderID, &order.Accrual, &order.Status, &order.CreatedAt)
		if err != nil {
			return []models.Order{}, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		logger.Log.Errorf("failed to scan rows for all orders: %v", err)
		return []models.Order{}, err
	}

	return orders, nil
}
