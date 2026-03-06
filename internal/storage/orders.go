package storage

import (
	"context"

	"gofermart_/internal/models"

	"github.com/Masterminds/squirrel"
)

// CreateOrder создаёт новый заказ для пользователя.
// Если заказ с таким номером уже существует, возвращает:
// - ErrOrderExistsForUser — если заказ уже загружен этим пользователем
// - ErrOrderExistsForOther — если заказ загружен другим пользователем
func (p *Postgres) CreateOrder(ctx context.Context, o *models.Order) error {
	query, args, err := p.sb.
		Insert("orders").
		Columns("number", "user_id", "status").
		Values(o.Number, o.UserID, o.Status).
		Suffix("ON CONFLICT (number) DO NOTHING").
		ToSql()
	if err != nil {
		return err
	}

	tag, err := p.pool.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 1 {
		return nil
	}

	queryS, args, err := p.sb.
		Select("user_id").
		From("orders").
		Where(squirrel.Eq{"number": o.Number}).
		ToSql()
	if err != nil {
		return err
	}

	var existingUserID int
	err = p.pool.QueryRow(ctx, queryS, args...).Scan(&existingUserID)
	if err != nil {
		return err
	}

	if existingUserID == o.UserID {
		return ErrOrderExistsForUser
	}

	return ErrOrderExistsForOther
}

// GetUserOrders возвращает все заказы пользователя, отсортированные по времени загрузки (DESC)
func (p *Postgres) GetUserOrders(ctx context.Context, userID int) ([]models.Order, error) {
	query, args, err := p.sb.
		Select("number", "status", "accrual", "uploaded_at").
		From("orders").
		Where(squirrel.Eq{"user_id": userID}).
		OrderBy("uploaded_at DESC").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := p.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var o models.Order
		o.UserID = userID

		if err := rows.Scan(
			&o.Number,
			&o.Status,
			&o.Accrual,
			&o.UploadedAt,
		); err != nil {
			return nil, err
		}

		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

// ClaimOrders выбирает до limit заказов со статусом NEW или PROCESSING (по update_at) и помечает их PROCESSING.
// Используется для обработки заказов фоновым воркером.
func (p *Postgres) ClaimOrders(ctx context.Context, limit int) ([]models.Order, error) {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	subQuery := p.sb.
		Select("number").
		From("orders").
		Where(
			squirrel.Or{
				squirrel.Eq{"status": "NEW"},
				squirrel.Expr("status = 'PROCESSING' AND updated_at < now() - interval '2 second'"),
			},
		).
		OrderBy("updated_at").
		Suffix("FOR UPDATE SKIP LOCKED").
		Limit(uint64(limit))

	query, args, err := p.sb.
		Update("orders").
		Set("status", "PROCESSING").
		Set("updated_at", squirrel.Expr("now()")).
		Where("number IN (?)", subQuery).
		Suffix("RETURNING number, user_id, status, accrual, uploaded_at").
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		if err := rows.Scan(
			&o.Number,
			&o.UserID,
			&o.Status,
			&o.Accrual,
			&o.UploadedAt,
		); err != nil {
			return nil, err
		}
		orders = append(orders, o)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return orders, nil
}

// UpdateOrderStatus обновляет статус заказа и начисляет accrual, если заказ PROCESSED.
// Создаёт запись в user_balance, если её ещё нет, и добавляет accrual к current.
func (p *Postgres) UpdateOrderStatus(ctx context.Context, number string, status models.OrderStatus, accrual float64) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	selectSQL, args, err := p.sb.
		Select("user_id").
		From("orders").
		Where(squirrel.Eq{"number": number}).
		Suffix("FOR UPDATE").
		ToSql()
	if err != nil {
		return err
	}

	var userID int
	err = tx.QueryRow(ctx, selectSQL, args...).Scan(&userID)
	if err != nil {
		return err
	}

	updateOrderSQL, args, err := p.sb.
		Update("orders").
		Set("status", status).
		Set("accrual", accrual).
		Where(squirrel.Eq{"number": number}).
		ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, updateOrderSQL, args...)
	if err != nil {
		return err
	}

	if status == models.OrderProcessed && accrual > 0 {
		insertBalanceSQL, args, err := p.sb.
			Insert("user_balance").
			Columns("user_id", "current", "withdrawn").
			Values(userID, 0, 0).
			Suffix("ON CONFLICT (user_id) DO NOTHING").
			ToSql()
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, insertBalanceSQL, args...)
		if err != nil {
			return err
		}

		updateBalanceSQL, args, err := p.sb.
			Update("user_balance").
			Set("current", squirrel.Expr("current + ?", accrual)).
			Where(squirrel.Eq{"user_id": userID}).
			ToSql()
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, updateBalanceSQL, args...)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
