package user_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

var (
	validID    = ulid.GenerateStaticULID("valid")
	validName  = "sato taro"
	validEmail = "sato@example.com"
)

func TestNew(t *testing.T) {
	tests := []struct {
		caseName string
		id       string
		name     string
		email    string
		wantErr  error
	}{
		{
			caseName: "happy path: 有効なUserエンティティを作成",
			id:       validID,
			name:     validName,
			email:    validEmail,
			wantErr:  nil,
		},
		{
			caseName: "edge case: 無効なIDを指定",
			id:       "invalid",
			name:     validName,
			email:    validEmail,
			wantErr:  user.ErrInvalidUserID,
		},
		{
			caseName: "edge case: 名前が最小文字数 -1",
			id:       validID,
			name:     strings.Repeat("a", user.NameMinLength-1),
			email:    validEmail,
			wantErr:  user.ErrInvalidUserName,
		},
		{
			caseName: "happy path: 名前が最小文字数",
			id:       validID,
			name:     strings.Repeat("a", user.NameMinLength),
			email:    validEmail,
			wantErr:  nil,
		},
		{
			caseName: "happy path: 名前が最大文字数",
			id:       validID,
			name:     strings.Repeat("a", user.NameMaxLength),
			email:    validEmail,
			wantErr:  nil,
		},
		{
			caseName: "edge case: 名前が最大文字数 +1",
			id:       validID,
			name:     strings.Repeat("a", user.NameMaxLength+1),
			email:    validEmail,
			wantErr:  user.ErrInvalidUserName,
		},
		{
			caseName: "edge case: 無効なEmailを指定",
			id:       validID,
			name:     validName,
			email:    "invalid",
			wantErr:  user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			u, err := user.New(tt.id, tt.name, tt.email)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.id, u.ID())
				assert.Equal(t, tt.name, u.Name())
				assert.Equal(t, tt.email, u.Email())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, u)
			}
		})
	}
}

func TestChangeName(t *testing.T) {
	tests := []struct {
		caseName string
		newName  string
		wantErr  error
	}{
		{
			caseName: "happy path: 有効な新しい名前に変更",
			newName:  "yamada hanako",
			wantErr:  nil,
		},
		{
			caseName: "edge case: 無効な名前を指定した時、名前は変更されない",
			newName:  strings.Repeat("a", user.NameMinLength-1),
			wantErr:  user.ErrInvalidUserName,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			u, _ := user.New(validID, validName, validEmail)
			err := u.ChangeName(tt.newName)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, u.Name())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, validName, u.Name())
			}
		})
	}
}

func TestChangeEmail(t *testing.T) {
	tests := []struct {
		caseName string
		newEmail string
		wantErr  error
	}{
		{
			caseName: "happy path: 有効なEmailに変更",
			newEmail: "yamada@example.com",
			wantErr:  nil,
		},
		{
			caseName: "edge case: 無効なEmailを指定した時、Emailは変更されない",
			newEmail: "invalid-email",
			wantErr:  user.ErrInvalidEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			u, _ := user.New(validID, validName, validEmail)
			err := u.ChangeEmail(tt.newEmail)

			if tt.wantErr == nil {
				assert.NoError(t, err)
				assert.Equal(t, tt.newEmail, u.Email())
			} else {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, validEmail, u.Email())
			}
		})
	}
}
