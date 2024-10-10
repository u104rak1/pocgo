package authentication_domain

import (
	user_domain "github.com/ucho456job/pocgo/internal/domain/user"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
)

type Authentication struct {
	userID       string
	passwordHash string
}

func New(userID, password string) (*Authentication, error) {
	if err := user_domain.ValidID(userID); err != nil {
		return nil, err
	}

	if err := ValidPassword(password); err != nil {
		return nil, err
	}
	passwordHash := passwordUtil.Encode(password)

	return &Authentication{
		userID:       userID,
		passwordHash: passwordHash,
	}, nil
}

func (p *Authentication) Compare(hash, password string) error {
	return passwordUtil.Compare(hash, password)
}
