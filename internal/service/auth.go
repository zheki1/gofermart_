package service

import (
	"context"
	"errors"

	"gofermart_/internal/auth"
	"gofermart_/internal/models"
	"gofermart_/internal/storage"
)

type AuthService interface {
	Register(ctx context.Context, c models.Credentials) (string, error)
	Login(ctx context.Context, c models.Credentials) (string, error)
}

type authService struct {
	repo storage.Repository
}

func NewAuthService(r storage.Repository) AuthService {
	return &authService{repo: r}
}

func (s *authService) Register(ctx context.Context, c models.Credentials) (string, error) {

	if c.Login == "" || c.Password == "" {
		return "", ErrInvalidCredentials
	}

	hash, err := auth.HashPassword(c.Password)
	if err != nil {
		return "", err
	}

	u := &models.User{
		Login:        c.Login,
		PasswordHash: hash,
	}

	err = s.repo.CreateUser(ctx, u)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return "", ErrUserExists
		}
		return "", err
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Login(ctx context.Context, c models.Credentials) (string, error) {

	if c.Login == "" || c.Password == "" {
		return "", ErrInvalidCredentials
	}

	u, err := s.repo.GetUserByLogin(ctx, c.Login)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := auth.CheckPassword(u.PasswordHash, c.Password); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := auth.GenerateToken(u.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}
