package storage

import (
	"context"

	"gofermart_/internal/models"
)

// CreateOrder создаёт новый заказ для пользователя.
// Если заказ с таким номером уже существует, возвращает:
// - ErrOrderExistsForUser — если заказ уже загружен этим пользователем
// - ErrOrderExistsForOther — если заказ загружен другим пользователем
func (p *Postgres) CreateOrder(ctx context.Context, o *models.Order) error {
	tag, err := p.pool.Exec(ctx,
		`INSERT INTO orders(number,user_id,status)
		 VALUES($1,$2,$3)
		 ON CONFLICT (number) DO NOTHING`,
		o.Number, o.UserID, o.Status,
	)
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 1 {
		return nil
	}

	var existingUserID int
	err = p.pool.QueryRow(ctx,
		`SELECT user_id FROM orders WHERE number=$1`,
		o.Number,
	).Scan(&existingUserID)
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
	rows, err := p.pool.Query(ctx,
		`SELECT number, status, accrual, uploaded_at
		 FROM orders
		 WHERE user_id=$1
		 ORDER BY uploaded_at DESC`,
		userID,
	)
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

	rows, err := tx.Query(ctx, `
		UPDATE orders
		SET status='PROCESSING', updated_at = now()
		WHERE number IN (
			SELECT number
			FROM orders
			WHERE status='NEW' OR (status='PROCESSING' AND updated_at < now() - interval '15 second')
			ORDER BY updated_at
			FOR UPDATE SKIP LOCKED
			LIMIT $1
		)
		RETURNING number, user_id, status, accrual, uploaded_at
	`, limit)
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

	// Получаем user_id заказа
	var userID int
	err = tx.QueryRow(ctx, `SELECT user_id FROM orders WHERE number=$1 FOR UPDATE`, number).Scan(&userID)
	if err != nil {
		return err
	}

	// Обновляем заказ
	_, err = tx.Exec(ctx,
		`UPDATE orders SET status=$1, accrual=$2 WHERE number=$3`,
		status, accrual, number,
	)
	if err != nil {
		return err
	}

	// Если заказ PROCESSED — обновляем баланс
	if status == models.OrderProcessed && accrual > 0 {
		// Создаём строку для пользователя, если ещё нет
		_, err = tx.Exec(ctx,
			`INSERT INTO user_balance(user_id, current, withdrawn)
			 VALUES($1, 0, 0)
			 ON CONFLICT (user_id) DO NOTHING`,
			userID,
		)
		if err != nil {
			return err
		}

		// Добавляем accrual к current
		_, err = tx.Exec(ctx,
			`UPDATE user_balance
			 SET current = current + $1
			 WHERE user_id=$2`,
			accrual, userID,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
