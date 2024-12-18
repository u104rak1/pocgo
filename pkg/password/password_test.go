package password_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/pkg/password"
	"golang.org/x/crypto/bcrypt"
)

func TestEncode(t *testing.T) {
	t.Run("パスワードをハッシュ化できる", func(t *testing.T) {
		_, err := password.Encode("ValidPassword123")
		assert.NoError(t, err)
	})
	// GenerateFromPassword関数でエラーを発生させることが難しいため、エラーケースのテストは省略します。
}

func TestCompare(t *testing.T) {
	passwordHash, _ := password.Encode("ValidPassword123")
	tests := []struct {
		caseName string
		password string
		hash     string
		errMsg   string
	}{
		{
			caseName: "パスワードが一致する場合は検証に成功する",
			password: "ValidPassword123",
			hash:     passwordHash,
			errMsg:   "",
		},
		{
			caseName: "パスワードが一致しない場合は検証に失敗する",
			password: "DifferentPassword456",
			hash:     passwordHash,
			errMsg:   bcrypt.ErrMismatchedHashAndPassword.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := password.Compare(tt.hash, tt.password)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
