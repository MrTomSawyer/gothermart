package workers

import (
	"encoding/json"
	"fmt"
	"github.com/MrTomSawyer/loyalty-system/internal/app/interfaces"
	"github.com/MrTomSawyer/loyalty-system/internal/app/logger"
	"github.com/MrTomSawyer/loyalty-system/internal/app/models"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

func HandleOrders(orderRep interfaces.OrderRepository, maxWorkers int, tickerPeriod string, accrualHost string, retryInterval string) {
	workerPool := make(chan struct{}, maxWorkers)

	go func() {
		tPeriod, err := strconv.Atoi(tickerPeriod)
		if err != nil {
			logger.Log.Errorf("Invalid ticker period format: %v", err)
			return
		}

		ticker := time.NewTicker(time.Duration(tPeriod) * time.Second)
		for range ticker.C {
			unhandledOrderIDs, err := orderRep.GetUnhandledOrders()
			if err != nil {
				logger.Log.Errorf("Failed to get all ")
			}
			for _, orderID := range unhandledOrderIDs {
				logger.Log.Infof("Procceeding order №%s", orderID)
				workerPool <- struct{}{}

				go processOrder(orderID, workerPool, orderRep, accrualHost, retryInterval)
			}
		}
	}()
}

func processOrder(orderID string, workerPool chan struct{}, orderRep interfaces.OrderRepository, accrualHost string, retryInterval string) {
	var m sync.Mutex
	m.Lock()
	defer m.Unlock()
	defer func() {
		<-workerPool
	}()

	userID, err := orderRep.GetOrderAndUserIDs(orderID)
	if err != nil {
		logger.Log.Errorf("no order with id=%s found: %v", orderID, err)
		return
	}

	order := models.Order{}
	order.OrderID = orderID
	getAccrual(&order, accrualHost, retryInterval)
	if order.Status == "" {
		logger.Log.Infof("No accrual received for order №%s", orderID)
		return
	}

	err = orderRep.UpdateOrderAccrual(order, orderID, userID)
	if err != nil {
		logger.Log.Errorf("failed to update order: %s", orderID)
		return
	}
}

func getAccrual(order *models.Order, accrualHost string, interval string) {
	url := fmt.Sprintf("%s/api/orders/%s", accrualHost, order.OrderID)
	logger.Log.Infof("Accrual url: %s", url)
	maxRetries := 5

	var res *http.Response
	var err error

	retryInt, err := strconv.Atoi(interval)
	if err != nil {
		logger.Log.Errorf("Invalid retry period format: %v", err)
	}
	retryInterval := time.Duration(retryInt) * time.Second

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
