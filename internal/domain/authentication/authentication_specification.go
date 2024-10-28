package authentication

import "errors"

const (
	PasswordMinLength = 8
	PasswordMaxLength = 20
)

var (
	ErrInvalidID                   = errors.New("authentication id must be a valid ULID")
	ErrPasswordInvalidLength       = errors.New("password must be between 8 and 20 characters")
	ErrAuthenticationAlreadyExists = errors.New("authentication already exists")
	ErrUnexpectedSigningMethod     = errors.New("unexpected signing method")
	ErrInvalidAccessToken          = errors.New("invalid access token")
	ErrAuthenticationFailed        = errors.New("email or password is incorrect")
	ErrAuthenticationNotFound      = errors.New("authentication not found")
	ErrUnmatchedPassword           = errors.New("passwords do not match")
)

func validPassword(password string) error {
	if len(password) < PasswordMinLength || len(password) > PasswordMaxLength {
		return ErrPasswordInvalidLength
	}
	return nil
}
