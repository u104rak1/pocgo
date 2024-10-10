package user

import "errors"

var (
	ErrInvalidUserID          = errors.New("invalid user id")
	ErrInvalidUserName        = errors.New("user name must be between 1 and 20 characters")
	ErrInvalidEmail           = errors.New("invalid email")
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
)
