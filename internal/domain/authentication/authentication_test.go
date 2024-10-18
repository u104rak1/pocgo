package authentication_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/domain/authentication"
	userDomain "github.com/ucho456job/pocgo/internal/domain/user"
	passwordUtil "github.com/ucho456job/pocgo/pkg/password"
	"github.com/ucho456job/pocgo/pkg/ulid"
	"golang.org/x/crypto/bcrypt"
)

var (
	validUserID   = ulid.GenerateStaticULID("user")
	validPassword = "password"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name     string
		userID   string
		password string
		wantErr  error
	}{
		{
			name:     "Happy path: 有効なAuthenticationを作成",
			userID:   validUserID,
			password: validPassword,
			wantErr:  nil,
		},
		{
			name:     "Edge case: 無効なUserIDを指定するとエラー",
			userID:   "invalid",
			password: validPassword,
			wantErr:  userDomain.ErrInvalidUserID,
		},
		{
			name:     "Edge case: パスワードが7文字だとエラー",
			userID:   validUserID,
			password: "1234567",
			wantErr:  authentication.ErrPasswordInvalidLength,
		},
		{
			name:     "Happy path: パスワードが8文字なら成功",
			userID:   validUserID,
			password: "12345678",
			wantErr:  nil,
		},
		{
			name:     "Happy path: パスワードが20文字なら成功",
			userID:   validUserID,
			password: "12345678901234567890",
			wantErr:  nil,
		},
		{
			name:     "Edge case: パスワードが21文字だとエラー",
			userID:   validUserID,
			password: "123456789012345678901",
			wantErr:  authentication.ErrPasswordInvalidLength,
		},
		// Encode関数のエラーを強制するのは難しい為、エラーのテストは省略
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth, err := authentication.New(tt.userID, tt.password)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
				assert.Nil(t, auth)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.userID, auth.UserID())
				assert.NoError(t, passwordUtil.Compare(auth.PasswordHash(), tt.password))
			}
		})
	}
}

func TestReconstruct(t *testing.T) {
	t.Run("Happy path: 有効なAuthenticationエンティティを再構築", func(t *testing.T) {
		encodedPassword, _ := passwordUtil.Encode(validPassword)
		auth, err := authentication.Reconstruct(validUserID, encodedPassword)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, validUserID, auth.UserID())
		assert.Equal(t, encodedPassword, auth.PasswordHash())
	})
}

func TestCompare(t *testing.T) {
	tests := []struct {
		name        string
		newPassword string
		wantErr     error
	}{
		{
			name:        "Happy path: パスワードが一致する場合、エラーが発生しない",
			newPassword: validPassword,
			wantErr:     nil,
		},
		{
			name:        "Edge case: パスワードが異なる場合、エラーが発生する",
			newPassword: "deffirentPassword",
			wantErr:     bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			auth, _ := authentication.New(validUserID, validPassword)
			err := auth.Compare(tt.newPassword)

			if tt.wantErr != nil {
				assert.ErrorIs(t, err, tt.wantErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
