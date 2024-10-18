package user_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/user"
	"github.com/ucho456job/pocgo/pkg/ulid"
)

var (
	validID    = ulid.GenerateStaticULID("user")
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
			caseName: "edge case: 無効なIDを指定するとエラー",
			id:       "invalid",
			name:     validName,
			email:    validEmail,
			wantErr:  user.ErrInvalidUserID,
		},
		{
			caseName: "edge case: 名前が0文字だとエラー",
			id:       validID,
			name:     strings.Repeat("a", user.NameMinLength-1),
			email:    validEmail,
			wantErr:  user.ErrInvalidUserName,
		},
		{
			caseName: "happy path: 名前が1文字なら成功",
			id:       validID,
			name:     strings.Repeat("a", user.NameMinLength),
			email:    validEmail,
			wantErr:  nil,
		},
		{
			caseName: "happy path: 名前が20文字なら成功",
			id:       validID,
			name:     strings.Repeat("a", user.NameMaxLength),
			email:    validEmail,
			wantErr:  nil,
		},
		{
			caseName: "edge case: 名前が21文字だとエラー",
			id:       validID,
			name:     strings.Repeat("a", user.NameMaxLength+1),
			email:    validEmail,
			wantErr:  user.ErrInvalidUserName,
		},
		{
			caseName: "edge case: 無効なEmailを指定するとエラー",
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
			caseName: "happy path: 有効な新しい名前に変更できる",
			newName:  "yamada hanako",
			wantErr:  nil,
		},
		{
			caseName: "edge case: 無効な名前を指定した時、名前を変更できない",
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
			caseName: "happy path: 有効なEmailに変更できる",
			newEmail: "yamada@example.com",
			wantErr:  nil,
		},
		{
			caseName: "edge case: 無効なEmailを指定した時、Emailを変更できない",
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
