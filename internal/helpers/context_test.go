package helpers

import (
	"context"
	"errors"
	"testing"

	"gofermart_/internal/http/middleware"
)

func TestGetUserID(t *testing.T) {

	tests := []struct {
		name    string
		ctx     context.Context
		wantID  int
		wantErr error
	}{
		{
			name: "valid user id",
			ctx: context.WithValue(
				context.Background(),
				middleware.UserIDKey,
				123,
			),
			wantID:  123,
			wantErr: nil,
		},
		{
			name:    "missing user id",
			ctx:     context.Background(),
			wantID:  0,
			wantErr: ErrUnauthorized,
		},
		{
			name: "wrong type",
			ctx: context.WithValue(
				context.Background(),
				middleware.UserIDKey,
				"123",
			),
			wantID:  0,
			wantErr: ErrUnauthorized,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			id, err := GetUserID(tt.ctx)

			if id != tt.wantID {
				t.Fatalf("expected id %d, got %d", tt.wantID, id)
			}

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
		})
	}
}
