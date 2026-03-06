package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gofermart_/internal/accrual"
	"gofermart_/internal/auth"
	"gofermart_/internal/config"
	httptransport "gofermart_/internal/http"
	"gofermart_/internal/logger"
	"gofermart_/internal/storage"
	"gofermart_/internal/worker"
)

// App представляет всё приложение Gofermart: конфигурацию, сервер, БД и воркера.
type App struct {
	cfg    *config.Config
	server *http.Server
	db     *storage.Postgres
	worker *worker.Worker
}

// New создаёт новый экземпляр приложения.
//
// 1. Инициализирует логгер.
// 2. Загружает конфигурацию из флагов и env.
// 3. Создаёт соединение с PostgreSQL и запускает миграции.
// 4. Создаёт клиента начислений (accrual) и воркера.
// 5. Настраивает HTTP-роутер и сервер.
//
// Возвращает указатель на App и ошибку при любой проблеме инициализации.
func New() (*App, error) {
	logger.Init("gofermart.log", logger.DEBUG)
	logger.Log.Info("Logger initialized")

	cfg := config.MustLoadWithFlags(nil)
	auth.Init(cfg.JWTSecret)

	db, err := storage.NewPostgres(cfg.DatabaseURI)
	if err != nil {
		return nil, fmt.Errorf("db init: %w", err)
	}

	if err := storage.RunMigrations(cfg.DatabaseURI); err != nil {
		return nil, fmt.Errorf("migrations: %w", err)
	}

	accrualClient := accrual.New(cfg.AccrualSystemAddress)
	w := worker.New(db, accrualClient)

	container := httptransport.NewContainer(db)
	router := httptransport.NewRouter(container)

	srv := &http.Server{
		Addr:         cfg.RunAddress,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		cfg:    cfg,
		server: srv,
		db:     db,
		worker: w,
	}, nil
}

// Run запускает HTTP-сервер и воркер.
//
// 1. Создаёт контекст, который реагирует на SIGINT/SIGTERM.
// 2. Запускает воркер в отдельной горутине.
// 3. Запускает сервер в отдельной горутине.
// 4. Ожидает сигнал завершения, затем корректно завершает сервер и БД.
//
// Возвращает ошибку, если завершение приложения прошло с ошибкой.
func (a *App) Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	go a.worker.Start(ctx)

	go func() {
		logger.Log.Info("Server started on", a.cfg.RunAddress)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("ListenAndServe error:", err)
			stop()
		}
	}()

	<-ctx.Done()
	logger.Log.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Error("Server shutdown error:", err)
	}

	if err := a.db.Close(shutdownCtx); err != nil {
		logger.Log.Error("DB close error:", err)
	}

	logger.Log.Info("Application stopped gracefully")

	return nil
}
