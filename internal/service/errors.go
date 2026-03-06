package service

import "errors"

var (
	ErrInvalidSum          = errors.New("invalid sum")
	ErrInvalidOrderNumber  = errors.New("invalid order number")
	ErrNotEnoughFunds      = errors.New("not enough funds")
	ErrInvalidCredentials  = errors.New("invalid login or password")
	ErrUserExists          = errors.New("login already exists")
	ErrOrderExistsForUser  = errors.New("order already uploaded by this user")
	ErrOrderExistsForOther = errors.New("order uploaded by another user")
)
