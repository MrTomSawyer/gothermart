package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MrTomSawyer/loyalty-system/internal/app/entity"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"io"
	"net/http"
	"sync"
	"time"
)

func HandleOrders(ctx context.Context, pool *pgxpool.Pool, orderCh chan string, maxWorkers int, accrualHost string) {
	var m sync.Mutex
	workerPool := make(chan struct{}, maxWorkers)

	for orderID := range orderCh {
		logger.Log.Infof("Procceeding order №%s", orderID)
		workerPool <- struct{}{}

		go func(orderID string) {
			m.Lock()
			defer m.Unlock()
			defer func() {
				<-workerPool
			}()

			// проверяем есть ли такой ордер в базе
			row := pool.QueryRow(ctx, "SELECT user_id, order_status from orders WHERE order_num=$1", orderID)
			var orderStatus string
			var userID int
			err := row.Scan(&userID)
			if err != nil {
				logger.Log.Errorf("no order with id=%d found: %v", orderID, err)
				return
			}
			if orderStatus == "INVALID" || orderStatus == "PROCESSED" {
				logger.Log.Errorf("order with id=%d has status of %s and can't be processed", orderID, orderStatus)
				return
			}

			order := models.Order{}
			order.OrderID = orderID
			getAccrual(&order, accrualHost)

			if order.Status == "" {
				logger.Log.Infof("No accrual received for order №%s", orderID)
				return
			}
			//апдейтим ордер и баланс юзера
			tx, err := pool.BeginTx(ctx, pgx.TxOptions{})
			if err != nil {
				logger.Log.Infof("failed to begin transaction: %v", err)
				return
			}
			//апдейтим ордер
			_, err = pool.Exec(ctx, "UPDATE orders SET order_status=$1, accrual=$2 WHERE order_num=$3", order.Status, order.Accrual, orderID)
			if err != nil {
				if err := tx.Rollback(ctx); err != nil {
					logger.Log.Errorf("failed to rollback transaction: %v", err)
				}
				logger.Log.Errorf("failed to update order=%d: %v", orderID, err)
			}

			//апдейтим находим в базе юзера и берем значение его баланса
			row = pool.QueryRow(ctx, "SELECT balance FROM users WHERE id=$1", userID)

			// восстанавливаем сущность юзера чтобы пополнить баланс
			var user entity.User
			err = row.Scan(&user.Balance)
			if err != nil {
				if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
					logger.Log.Errorf("failed to rollback transaction: %v", rollbackErr)
				}
				logger.Log.Errorf("failed to find user with id=%d : %v", userID, err)
			}

			// пополняем баланс
			_ = user.SetBalance(order.Accrual)

			// апдейтим юзера
			_, err = pool.Exec(ctx, "UPDATE users SET balance=$1 WHERE id=$2", user.Balance, userID)
			if err != nil {
				if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
					logger.Log.Errorf("failed to rollback transaction: %v", err)
				}
				logger.Log.Errorf("failed to upfate user with id=%d : %v", userID, err)
			}

			// коммитим транзакцию.
			if err := tx.Commit(ctx); err != nil {
				logger.Log.Errorf("failed to commit transaction: %v", err)
			}
			logger.Log.Infof("Transaction successfully commited")

		}(orderID)
	}
}

func getAccrual(order *models.Order, accrualHost string) {
	url := fmt.Sprintf("%s/api/orders/%s", accrualHost, order.OrderID)
	logger.Log.Infof("Accrual url: %s", url)
	maxRetries := 5
	retryInterval := 1 * time.Second

	var res *http.Response
	var err error

	for maxRetries > 0 {
		res, err = http.Get(url)
		if err != nil {
			logger.Log.Errorf("failed to get accrual info: %v", err)
			return
		}

		switch res.StatusCode {
		case http.StatusOK:
			body, err := io.ReadAll(res.Body)
			if err != nil {
				logger.Log.Infof("Error reading response")
				break
			}
			var accrual models.Accrual
			err = json.Unmarshal(body, &accrual)
			if err != nil {
				logger.Log.Infof("failed to unmarshal responce body: %v", err)
				logger.Log.Infof("Response body: %s", body)
				break
			}
			order.Status = accrual.Status
			order.Accrual = accrual.Accrual
			return
		case http.StatusNoContent:
			logger.Log.Errorf("no order with id=%s found in the accrual system", order.OrderID)
			return
		case http.StatusTooManyRequests:
			logger.Log.Errorf("too many requests. Retrying...")
			maxRetries--
			time.Sleep(retryInterval)
			continue
		case http.StatusInternalServerError:
			logger.Log.Infof("internal server error. Retrying...")
			maxRetries--
			time.Sleep(retryInterval)
			continue
		default:
			logger.Log.Infof("something weng wrong. response status: %s. Retrying...", res.Status)
			maxRetries--
			time.Sleep(retryInterval)
			continue
		}

		if res != nil {
			err = res.Body.Close()
			if err != nil {
				logger.Log.Errorf("error closing response body")
			}
		}
		return
	}
}
