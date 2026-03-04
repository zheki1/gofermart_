package accrual

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Client представляет HTTP-клиент для взаимодействия с начислениями.
type Client struct {
	baseURL string
	client  *http.Client
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
	return &Client{
		baseURL: baseURL,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
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
func (c *Client) GetOrder(number string) (*Response, int, time.Duration, error) {
	url := fmt.Sprintf("http://%s/api/orders/%s", c.baseURL, number)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, 0, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := 0 * time.Second
		if header := resp.Header.Get("Retry-After"); header != "" {
			header = strings.TrimSpace(header)
			if seconds, err := strconv.Atoi(header); err == nil {
				retryAfter = time.Duration(seconds) * time.Second
			}
		}
		return nil, resp.StatusCode, retryAfter, nil
	}

	if resp.StatusCode == http.StatusNoContent {
		return nil, resp.StatusCode, 0, nil
	}

	var result Response
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, resp.StatusCode, 0, err
	}

	return &result, resp.StatusCode, 0, nil
}
