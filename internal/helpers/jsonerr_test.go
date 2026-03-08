package helpers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteJSONError(t *testing.T) {
	// создаём ResponseRecorder для имитации http.ResponseWriter
	rr := httptest.NewRecorder()

	status := http.StatusBadRequest
	msg := "invalid request format"

	WriteJSONError(rr, msg, status)

	res := rr.Result()
	defer res.Body.Close()

	// проверяем статус код
	if res.StatusCode != status {
		t.Errorf("expected status %d, got %d", status, res.StatusCode)
	}

	// проверяем заголовок Content-Type
	if ct := res.Header.Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type 'application/json', got '%s'", ct)
	}

	// проверяем тело ответа
	var body map[string]string
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		t.Fatalf("failed to decode response body: %v", err)
	}

	if body["error"] != msg {
		t.Errorf("expected error message '%s', got '%s'", msg, body["error"])
	}
}
