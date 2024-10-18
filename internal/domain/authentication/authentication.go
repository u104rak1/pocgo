package authentication

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
	if err := validPassword(password); err != nil {
		return nil, err
	}
	passwordHash, err := passwordUtil.Encode(password)
	if err != nil {
		return nil, err
	}

	return newAuthentication(userID, passwordHash)
}

// DBからの再構築時は既にハッシュ値なのでパスワードのバリデーションを行わない
func Reconstruct(userID, passwordHash string) (*Authentication, error) {
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

func (a *Authentication) Compare(password string) error {
	return passwordUtil.Compare(a.passwordHash, password)
}
