package authentication

import (
	user_domain "github.com/ucho456job/pocgo/internal/domain/user"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
)

type Authentication struct {
	userID       string
	passwordHash string
}

func New(userID, password string) (*Authentication, error) {
	// Validate password when creating a new authentication
	if err := validPassword(password); err != nil {
		return nil, err
	}
	passwordHash, err := passwordUtil.Encode(password)
	if err != nil {
		return nil, err
	}

	return newAuthentication(userID, passwordHash)
}

func Reconstruct(userID, passwordHash string) (*Authentication, error) {
	// When reconstructing the authentication from the DB, the password is already encoded so there is no validation.
	return newAuthentication(userID, passwordHash)
}

func newAuthentication(userID, passwordHash string) (*Authentication, error) {
	if err := user_domain.ValidID(userID); err != nil {
		return nil, err
	}
	return &Authentication{
		userID:       userID,
		passwordHash: passwordHash,
	}, nil
}

func (a *Authentication) UserID() string {
	return a.userID
}

func (a *Authentication) PasswordHash() string {
	return a.passwordHash
}

func (a *Authentication) ComparePassword(password string) error {
	if err := passwordUtil.Compare(a.passwordHash, password); err != nil {
		return ErrUnmatchedPassword
	}
	return nil
}
