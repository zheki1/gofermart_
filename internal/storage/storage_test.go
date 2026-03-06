package storage

import (
	"context"
	"testing"
)

// Тест базовых ошибок
func TestStorageErrors(t *testing.T) {
	if ErrUserExists.Error() != "user already exists" {
		t.Error("ErrUserExists has wrong message")
	}
	if ErrUserNotFound.Error() != "user not found" {
		t.Error("ErrUserNotFound has wrong message")
	}
	if ErrNotEnoughFunds.Error() != "not enough funds" {
		t.Error("ErrNotEnoughFunds has wrong message")
	}
	if ErrOrderAlreadyWithdrawn.Error() != "order already withdrawn" {
		t.Error("ErrOrderAlreadyWithdrawn has wrong message")
	}
}

// Тест заглушки Postgres
func TestPostgres_Close(t *testing.T) {
	p := &Postgres{}
	if err := p.Close(context.Background()); err != nil {
		t.Errorf("Close should not return error, got %v", err)
	}
}

// RunMigrations не тестируем реально, тк требует БД
// но проверим что функция существует и не паникует при пустой строке DSN
func TestRunMigrations_InvalidDSN(t *testing.T) {
	err := RunMigrations("")
	if err == nil {
		t.Error("expected error for empty DSN")
	}
}
