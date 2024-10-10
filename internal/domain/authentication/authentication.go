package authentication_domain

import (
	user_domain "github.com/ucho456job/pocgo/internal/domain/user"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
)

type Authentication struct {
	userID       string
	passwordHash string
}

// 新規作成時はパスワードのバリデーションを行う
func New(userID, password string) (*Authentication, error) {
	if err := ValidPassword(password); err != nil {
		return nil, err
	}
	passwordHash, err := passwordUtil.Encode(password)
	if err != nil {
		return nil, err
	}

	return new(userID, passwordHash)
}

// DBからの再構築時は既にハッシュ値なのでパスワードのバリデーションを行わない
func Reconstruct(userID, passwordHash string) (*Authentication, error) {
	return new(userID, passwordHash)
}

func new(userID, passwordHash string) (*Authentication, error) {
	if err := user_domain.ValidID(userID); err != nil {
		return nil, err
	}
	return &Authentication{
		userID:       userID,
		passwordHash: passwordHash,
	}, nil
}

func (p *Authentication) UserID() string {
	return p.userID
}

func (p *Authentication) PasswordHash() string {
	return p.passwordHash
}

func (p *Authentication) Compare(hash, password string) error {
	return passwordUtil.Compare(hash, password)
}
