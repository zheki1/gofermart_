package worker

import (
	"context"
	"sync"
	"time"

	"gofermart_/internal/accrual"
	"gofermart_/internal/logger"
	"gofermart_/internal/models"
	"gofermart_/internal/storage"
)

// Worker — основной объект фонового воркера.
//
// repo — интерфейс работы с заказами и балансом.
// accrualClient — клиент внешней системы начислений.
// concurrency — число одновременно обрабатываемых заказов.
// interval — период опроса базы.
// wg — wait group для горутин.
type Worker struct {
	repo          storage.Repository
	accrualClient AccrualClient
	concurrency   int
	interval      time.Duration
	wg            sync.WaitGroup
}

type AccrualClient interface {
	GetOrder(number string) (*accrual.Response, int, time.Duration, error)
}

// New создаёт нового воркера с дефолтными параметрами.
//
// repo — интерфейс Repository для работы с заказами.
// accrualClient — клиент системы начислений.
func New(repo storage.Repository, accrualClient AccrualClient) *Worker {
	return &Worker{
		repo:          repo,
		accrualClient: accrualClient,
		concurrency:   5,               // число потоков
		interval:      5 * time.Second, // период тикера
	}
}

// Start запускает воркер в цикле.
// ctx используется для graceful shutdown.
//
// Внутри стартует ticker с интервалом w.interval, который вызывает
// метод process() для обработки заказов.
// handleOrder обрабатывает один заказ через accrualClient.
// - выполняет retry с backoff при ошибках соединения
// - учитывает код 429 и header Retry-After
// - обновляет статус заказа через UpdateOrderStatus
func (w *Worker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	logger.Log.Info("Worker started")

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Worker shutting down...")
			w.wg.Wait()
			logger.Log.Info("Worker stopped")
			return
		case <-ticker.C:
			w.process(ctx)
		}
	}
}

// process получает заказы через ClaimOrders и запускает для каждого
// go-функцию handleOrder.
func (w *Worker) process(ctx context.Context) {
	orders, err := w.repo.ClaimOrders(ctx, w.concurrency)
	if err != nil {
		logger.Log.Error("get pending orders:", err)
		return
	}

	for _, o := range orders {
		w.wg.Add(1)
		go func(order models.Order) {
			defer w.wg.Done()
			w.handleOrder(ctx, order)
		}(o)
	}
}

// handleOrder обрабатывает один заказ через accrualClient.
// - выполняет retry с backoff при ошибках соединения
// - учитывает код 429 и header Retry-After
// - обновляет статус заказа через UpdateOrderStatus
func (w *Worker) handleOrder(ctx context.Context, o models.Order) {
	backoff := time.Second
	maxBackoff := 30 * time.Second

	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		resp, code, retryAfter, err := w.accrualClient.GetOrder(o.Number)
		if err != nil {
			logger.Log.Error("accrual error:", "order", o.Number, "err", err, "backoff", backoff)

			wait := backoff
			if wait > maxBackoff {
				wait = maxBackoff
			}

			if !w.sleep(ctx, wait) {
				return
			}

			backoff *= 2
			continue
		}

		if code == 429 {
			wait := retryAfter
			if wait == 0 {
				wait = 5 * time.Second
			}
			logger.Log.Info("rate limited:", "order", o.Number, "wait", wait)

			if !w.sleep(ctx, wait) {
				return
			}
			continue
		}

		if resp == nil {
			logger.Log.Error("resp from accrual client is nil")
			return
		}

		var accrualValue float64
		if resp.Accrual != nil {
			accrualValue = *resp.Accrual
		}

		switch resp.Status {
		case "INVALID":
			w.repo.UpdateOrderStatus(ctx, o.Number, models.OrderInvalid, 0)
			logger.Log.Info("Order invalid:", o.Number)
		case "PROCESSING", "NEW":
			// уже PROCESSING, ничего не меняем
			logger.Log.Debug("Order still processing:", o.Number)
		case "PROCESSED":
			w.repo.UpdateOrderStatus(ctx, o.Number, models.OrderProcessed, accrualValue)
			logger.Log.Info("Order processed:", o.Number, "accrual:", accrualValue)
		}

		return
	}
}

// process получает заказы через ClaimOrders и запускает для каждого
// go-функцию handleOrder.
func (w *Worker) sleep(ctx context.Context, d time.Duration) bool {
	select {
	case <-time.After(d):
		return true
	case <-ctx.Done():
		return false
	}
}
