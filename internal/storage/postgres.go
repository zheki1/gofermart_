package storage

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

var _ Repository = (*Postgres)(nil)

// Postgres — реализация Repository для PostgreSQL.
type Postgres struct {
	pool *pgxpool.Pool
}

// NewPostgres создает новый Postgres и подключается к базе по URI.
// Возвращает ошибку при невозможности подключения или ping.
func NewPostgres(uri string) (*Postgres, error) {
	pool, err := pgxpool.New(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return &Postgres{pool: pool}, nil
}

// Close закрывает подключение к базе данных.
func (p *Postgres) Close(ctx context.Context) error {
	if p.pool != nil {
		p.pool.Close()
	}
	return nil
}
