package accrual

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_GetOrder(t *testing.T) {
	tests := []struct {
		name       string
		serverResp int
		respBody   interface{}
		wantStatus int
		wantRetry  time.Duration
		wantErr    bool
	}{
		{
			name:       "200 OK with data",
			serverResp: http.StatusOK,
			respBody: Response{
				Order:   "12345678903",
				Status:  "PROCESSED",
				Accrual: ptrFloat(100),
			},
			wantStatus: http.StatusOK,
			wantRetry:  0,
			wantErr:    false,
		},
		{
			name:       "204 No Content",
			serverResp: http.StatusNoContent,
			respBody:   nil,
			wantStatus: http.StatusNoContent,
			wantRetry:  0,
			wantErr:    false,
		},
		{
			name:       "429 Too Many Requests with Retry-After",
			serverResp: http.StatusTooManyRequests,
			respBody:   nil,
			wantStatus: http.StatusTooManyRequests,
			wantRetry:  3 * time.Second,
			wantErr:    false,
		},
		{
			name:       "invalid JSON",
			serverResp: http.StatusOK,
			respBody:   "{invalid json}",
			wantStatus: http.StatusOK,
			wantRetry:  0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.serverResp == http.StatusTooManyRequests {
					w.Header().Set("Retry-After", "3")
				}
				w.WriteHeader(tt.serverResp)
				switch v := tt.respBody.(type) {
				case Response:
					_ = json.NewEncoder(w).Encode(v)
				case string:
					w.Write([]byte(v))
				}
			})

			srv := httptest.NewServer(handler)
			defer srv.Close()

			client := New(srv.Listener.Addr().String())
			res, status, retry, err := client.GetOrder("12345678903")

			if status != tt.wantStatus {
				t.Errorf("status = %d, want %d", status, tt.wantStatus)
			}
			if retry != tt.wantRetry {
				t.Errorf("retry = %v, want %v", retry, tt.wantRetry)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.serverResp == http.StatusOK && !tt.wantErr && (res == nil || res.Order != "12345678903") {
				t.Errorf("response = %+v", res)
			}
		})
	}
}

func ptrFloat(f float64) *float64 { return &f }
