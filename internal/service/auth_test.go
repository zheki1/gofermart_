package service

import (
	"context"
	"testing"

	"gofermart_/internal/auth"
	"gofermart_/internal/models"

	"github.com/stretchr/testify/require"
)

func init() {
	auth.Init("test-secret")
}

func TestRegister_Success(t *testing.T) {
	s := NewAuthService(&mockRepo{})

	token, err := s.Register(context.Background(), models.Credentials{
		Login:    "user",
		Password: "password",
	})

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestRegister_InvalidCredentials(t *testing.T) {
	s := NewAuthService(&mockRepo{})

	_, err := s.Register(context.Background(), models.Credentials{})

	require.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_Success(t *testing.T) {
	s := NewAuthService(&mockRepo{})

	token, err := s.Login(context.Background(), models.Credentials{
		Login:    "user",
		Password: "password",
	})

	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestLogin_InvalidCredentials(t *testing.T) {
	s := NewAuthService(&mockRepo{})

	_, err := s.Login(context.Background(), models.Credentials{})

	require.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestLogin_WrongPassword(t *testing.T) {
	s := NewAuthService(&mockRepo{})

	_, err := s.Login(context.Background(), models.Credentials{
		Login:    "user",
		Password: "wrong",
	})

	require.ErrorIs(t, err, ErrInvalidCredentials)
}
