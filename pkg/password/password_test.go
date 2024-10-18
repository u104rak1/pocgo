package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/pkg/password"
	"golang.org/x/crypto/bcrypt"
)

func TestEncode(t *testing.T) {
	t.Run("Valid: 正常なパスワードをハッシュ化", func(t *testing.T) {
		_, err := password.Encode("ValidPassword123")
		assert.NoError(t, err)
	})
	// GenerateFromPassword関数のエラーを強制するのは難しい為、エラーのテストは省略
}

func TestCompare(t *testing.T) {
	passwordHash, _ := password.Encode("ValidPassword123")
	tests := []struct {
		caseName string
		password string
		hash     string
		wantErr  error
	}{
		{
			caseName: "Valid: 一致するパスワードで比較",
			password: "ValidPassword123",
			hash:     passwordHash,
			wantErr:  nil,
		},
		{
			caseName: "Invalid: 異なるパスワードを指定するとエラー",
			password: "DifferentPassword456",
			hash:     passwordHash,
			wantErr:  bcrypt.ErrMismatchedHashAndPassword,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := password.Compare(tt.hash, tt.password)

			if tt.wantErr != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
