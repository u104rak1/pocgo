package authentication_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	passwordUtil "github.com/u104rak1/pocgo/pkg/password"
)

func TestNew(t *testing.T) {
	var (
		userID   = idVO.NewUserIDForTest("user")
		password = "password"
	)

	tests := []struct {
		caseName string
		userID   idVO.UserID
		password string
		errMsg   string
	}{
		{
			caseName: "Positive: 認証情報を作成できる",
			userID:   userID,
			password: password,
			errMsg:   "",
		},
		{
			caseName: "Negative: 7文字のパスワードの場合はエラーが返る",
			userID:   userID,
			password: "1234567",
			errMsg:   "password must be between 8 and 20 characters",
		},
		{
			caseName: "Positive: 8文字のパスワードの場合は認証情報を作成できる",
			userID:   userID,
			password: "12345678",
			errMsg:   "",
		},
		{
			caseName: "Positive: 20文字のパスワードの場合は認証情報を作成できる",
			userID:   userID,
			password: "12345678901234567890",
			errMsg:   "",
		},
		{
			caseName: "Negative: 21文字のパスワードの場合はエラーが返る",
			userID:   userID,
			password: "123456789012345678901",
			errMsg:   "password must be between 8 and 20 characters",
		},
		// Password.Encode関数を強制的にエラーにすることが難しい為、このエラーパターンはテストしない
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			auth, err := authDomain.New(tt.userID, tt.password)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
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
	var (
		userID   = idVO.NewUserIDForTest("user").String()
		password = "password"
	)
	t.Run("Positive: 認証情報を再構築できる", func(t *testing.T) {
		encodedPassword, _ := passwordUtil.Encode(password)
		auth, err := authDomain.Reconstruct(userID, encodedPassword)

		assert.NoError(t, err)
		assert.NotNil(t, auth)
		assert.Equal(t, userID, auth.UserIDString())
		assert.Equal(t, encodedPassword, auth.PasswordHash())
	})
}

func TestComparePassword(t *testing.T) {
	var (
		userID   = idVO.NewUserIDForTest("user")
		password = "password"
	)

	tests := []struct {
		caseName    string
		newPassword string
		errMsg      string
	}{
		{
			caseName:    "Positive: パスワードが一致しない場合はエラーが返る",
			newPassword: password,
			errMsg:      "",
		},
		{
			caseName:    "Negative: パスワードが一致しない場合はエラーが返る",
			newPassword: "deffirentPassword",
			errMsg:      "passwords do not match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			auth, _ := authDomain.New(userID, password)
			err := auth.ComparePassword(tt.newPassword)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
