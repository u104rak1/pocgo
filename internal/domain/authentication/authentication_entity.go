package password_domain

import (
	user_domain "github.com/ucho456job/pocgo/internal/domain/user"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
)

type Authentication struct {
	id           string
	userID       string
	passwordHash string
}

func New(id, userID, password string) (*Authentication, error) {
	var err error
	if err = IsValidID(id); err != nil {
		return nil, err
	}

	if err = user_domain.IsValidID(userID); err != nil {
		return nil, err
	}

	if err = isValidPassword(password); err != nil {
		return nil, err
	}
	passwordHash := passwordUtil.Encode(password)

	return &Authentication{
		id:           id,
		userID:       userID,
		passwordHash: passwordHash,
	}, nil
}

func (p *Authentication) Compare(hash, password string) error {
	return passwordUtil.Compare(hash, password)
}
