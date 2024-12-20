package authentication

import (
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
)

type Authentication struct {
	userID       idVO.UserID
	passwordHash string
}

// 認証エンティティを作成します。新規で作成するのでパスワードの検証とハッシュ化を行います。
func New(userID idVO.UserID, password string) (*Authentication, error) {
	if err := validPassword(password); err != nil {
		return nil, err
	}
	passwordHash, err := passwordUtil.Encode(password)
	if err != nil {
		return nil, err
	}

	return newAuthentication(userID, passwordHash)
}

// データベースから認証を再構築します。パスワードは既にエンコードされているため、検証は行われません。
func Reconstruct(userID, passwordHash string) (*Authentication, error) {
	uID, err := idVO.UserIDFromString(userID)
	if err != nil {
		return nil, err
	}
	return newAuthentication(uID, passwordHash)
}

func newAuthentication(userID idVO.UserID, passwordHash string) (*Authentication, error) {
	return &Authentication{
		userID:       userID,
		passwordHash: passwordHash,
	}, nil
}

func (a *Authentication) UserID() idVO.UserID {
	return a.userID
}

func (a *Authentication) UserIDString() string {
	return a.userID.String()
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
