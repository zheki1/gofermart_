package accrual

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
)

// Client представляет HTTP-клиент для взаимодействия с начислениями.
type Client struct {
	baseURL string
	client  *resty.Client
}

// Response описывает ответ системы начислений для одного заказа.
type Response struct {
	Order   string   `json:"order"`             // номер заказа
	Status  string   `json:"status"`            // статус заказа (NEW, PROCESSING, PROCESSED, INVALID)
	Accrual *float64 `json:"accrual,omitempty"` // сумма начисления, если заказ PROCESSED
}

// New создаёт новый экземпляр Client.
//
// baseURL — адрес сервиса начислений (host:port)
// Возвращает готовый клиент с таймаутом 5 секунд.
func New(baseURL string) *Client {
	c := resty.New()
	c.SetTimeout(5 * time.Second)

	return &Client{
		baseURL: baseURL,
		client:  c,
	}
}

// GetOrder запрашивает информацию о заказе в системе начислений.
//
// number — номер заказа для запроса.
//
// Возвращает:
// - указатель на Response с данными заказа (nil, если заказ не найден или нет данных)
// - HTTP статус код
// - время, которое нужно ждать до повторного запроса (если 429 Too Many Requests)
// - ошибку (если запрос или декодирование JSON завершились неудачно)
func (c *Client) GetOrder(number string) (*Response, error) {
	var (
		result     *Response
		statusCode int
		retryAfter time.Duration
		lastErr    error
	)

	url := fmt.Sprintf("http://%s/api/orders/%s", c.baseURL, number)

	// retry-go для повторов при 429 или временных ошибках
	err := retry.Do(
		func() error {
			resp, err := c.client.R().Get(url)
			if err != nil {
				lastErr = err
				return err
			}

			statusCode = resp.StatusCode()

			if statusCode == 429 {
				// читаем Retry-After
				retryAfter = 0
				if header := resp.Header().Get("Retry-After"); header != "" {
					if seconds, err := strconv.Atoi(header); err == nil {
						retryAfter = time.Duration(seconds) * time.Second
					}
				}
				lastErr = fmt.Errorf("too many requests, retry after %v", retryAfter)
				return lastErr // retry
			}

			if statusCode == 204 {
				result = nil
				return nil // success, но пустой результат
			}

			var r Response
			if err := json.Unmarshal(resp.Body(), &r); err != nil {
				lastErr = err
				return err
			}

			result = &r
			return nil
		},
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			if retryAfter > 0 {
				d := retryAfter
				retryAfter = 0
				return d
			}
			// передаем err как второй аргумент
			return retry.BackOffDelay(n, err, config)
		}),
		retry.Attempts(3), // максимум 3 попытки
	)

	return result, err
}
