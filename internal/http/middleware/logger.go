package middleware

import (
	"bytes"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"gofermart_/internal/logger"
)

const maxErrorBodySize = 1024

// responseWriter оборачивает http.ResponseWriter для записи статуса, размера и тела ответа
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
	body   bytes.Buffer
}

// WriteHeader сохраняет статус ответа
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write сохраняет тело ответа для ошибок
func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.status == 0 {
		rw.status = http.StatusOK
	}

	if rw.status >= http.StatusBadRequest {
		if rw.body.Len() < maxErrorBodySize {
			remaining := maxErrorBodySize - rw.body.Len()
			if len(b) > remaining {
				rw.body.Write(b[:remaining])
			} else {
				rw.body.Write(b)
			}
		}
	}

	n, err := rw.ResponseWriter.Write(b)
	rw.size += n
	return n, err
}

// Logging логирует все запросы, ошибки и паники.
// Использует logger.Log для записи информации.
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rw := &responseWriter{
			ResponseWriter: w,
		}

		defer func() {
			if rec := recover(); rec != nil {
				logger.Log.Error(
					"PANIC",
					"method=", r.Method,
					"path=", r.URL.Path,
					"remote=", r.RemoteAddr,
					"panic=", rec,
					"stack=", string(debug.Stack()),
				)

				http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}

			if rw.status == 0 {
				rw.status = http.StatusOK
			}

			duration := time.Since(start)

			args := []interface{}{
				"method=", r.Method,
				"path=", r.URL.Path,
				"status=", rw.status,
				"size=", rw.size,
				"duration=", duration,
				"remote=", r.RemoteAddr,
			}

			switch {
			case rw.status >= 500:
				args = append(args, "error=", strings.TrimSpace(rw.body.String()))
				logger.Log.Error(args...)

			case rw.status >= 400:
				args = append(args, "error=", strings.TrimSpace(rw.body.String()))
				logger.Log.Warn(args...)

			default:
				logger.Log.Info(args...)
			}
		}()

		next.ServeHTTP(rw, r)
	})
}
