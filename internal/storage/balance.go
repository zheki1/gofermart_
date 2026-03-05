package storage

import (
	"context"
	"errors"
	"gofermart_/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (p *Postgres) Withdraw(ctx context.Context, userID int, order string, sum float64) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO withdrawals(user_id, order_number, sum, processed_at)
		VALUES($1,$2,$3,NOW())
	`, userID, order, sum)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return ErrOrderAlreadyWithdrawn
		}
		return err
	}

	tag, err := tx.Exec(ctx, `
		UPDATE user_balance
		SET current = current - $1,
		    withdrawn = withdrawn + $1
		WHERE user_id = $2
		  AND current >= $1
	`, sum, userID)

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

	err := p.pool.QueryRow(ctx, `
		SELECT current, withdrawn
		FROM user_balance
		WHERE user_id=$1
	`, userID).Scan(&current, &withdrawn)

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
	rows, err := p.pool.Query(ctx, `
		SELECT order_number, sum, processed_at
		FROM withdrawals
		WHERE user_id=$1
		ORDER BY processed_at DESC
	`, userID)
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
