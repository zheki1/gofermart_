package storage

import (
	"context"
	"errors"
	"gofermart_/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) Withdraw(ctx context.Context, userID int, order string, sum float64) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query, args, err := p.sb.
		Insert("withdrawals").
		Columns("user_id", "order_number", "sum", "processed_at").
		Values(userID, order, sum, squirrel.Expr("NOW()")).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return ErrOrderAlreadyWithdrawn
		}
		return err
	}

	queryUp, args, err := p.sb.
		Update("user_balance").
		Set("current", squirrel.Expr("current - ?", sum)).
		Set("withdrawn", squirrel.Expr("withdrawn + ?", sum)).
		Where(squirrel.Eq{"user_id": userID}).
		Where("current >= ?", sum).
		ToSql()
	if err != nil {
		return err
	}

	tag, err := tx.Exec(ctx, queryUp, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return ErrNotEnoughFunds
	}

	return tx.Commit(ctx)
}

func (p *Postgres) GetBalance(ctx context.Context, userID int) (*models.Balance, error) {
	var current float64
	var withdrawn float64

	query, args, err := p.sb.
		Select("current", "withdrawn").
		From("user_balance").
		Where(squirrel.Eq{"user_id": userID}).
		ToSql()
	if err != nil {
		return nil, err
	}

	err = p.pool.
		QueryRow(ctx, query, args...).
		Scan(&current, &withdrawn)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &models.Balance{
		Current:   current,
		Withdrawn: withdrawn,
	}, nil
}

func (p *Postgres) GetWithdrawals(ctx context.Context, userID int) ([]models.Withdrawal, error) {
	query, args, err := p.sb.
		Select("order_number", "sum", "processed_at").
		From("withdrawals").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("processed_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make([]models.Withdrawal, 0)

	for rows.Next() {
		var w models.Withdrawal
		if err := rows.Scan(&w.Order, &w.Sum, &w.ProcessedAt); err != nil {
			return nil, err
		}
		result = append(result, w)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
