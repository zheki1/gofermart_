package storage

import (
	"context"
	"errors"

	"gofermart_/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// CreateUser добавляет нового пользователя в базу данных.
// Если пользователь с таким логином уже существует, возвращает ErrUserExists.
// При успешном добавлении обновляет поля u.ID и u.CreatedAt.
func (p *Postgres) CreateUser(ctx context.Context, u *models.User) error {
	query, args, err := p.sb.
		Insert("users").
		Columns("login", "password_hash").
		Values(u.Login, u.PasswordHash).
		Suffix("RETURNING id, created_at").
		ToSql()
	if err != nil {
		return err
	}

	err = p.pool.QueryRow(ctx, query, args...).Scan(&u.ID, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return ErrUserExists
		}
		return err
	}
	return nil
}

// GetUserByLogin извлекает пользователя из базы по логину.
// Если пользователь не найден, возвращает ErrUserNotFound.
func (p *Postgres) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	u := &models.User{}

	query, args, err := p.sb.
		Select("id", "login", "password_hash", "created_at").
		From("users").
		Where(squirrel.Eq{"login": login}).
		ToSql()
	if err != nil {
		return nil, err
	}

	// Выполняем запрос
	err = p.pool.QueryRow(ctx, query, args...).Scan(
		&u.ID, &u.Login, &u.PasswordHash, &u.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return u, nil
}
