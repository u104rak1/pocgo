package user_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/pkg/ulid"
)

func TestNew(t *testing.T) {
	var (
		name  = "sato taro"
		email = "sato@example.com"
	)

	tests := []struct {
		caseName string
		id       string
		name     string
		email    string
		errMsg   string
	}{
		{
			caseName: "Positive: ユーザーが作成できる",
			name:     name,
			email:    email,
			errMsg:   "",
		},
		{
			caseName: "Negative: 名前が3文字未満の場合はエラーが返る",
			name:     strings.Repeat("a", 2),
			email:    email,
			errMsg:   "user name must be between 3 and 20 characters",
		},
		{
			caseName: "Positive: 名前が3文字の場合は、ユーザーが作成できる",
			name:     strings.Repeat("a", 3),
			email:    email,
			errMsg:   "",
		},
		{
			caseName: "Positive: 名前が20文字の場合は、ユーザーが作成できる",
			name:     strings.Repeat("a", 20),
			email:    email,
			errMsg:   "",
		},
		{
			caseName: "Negative: 名前が21文字の場合はエラーが返る",
			name:     strings.Repeat("a", 21),
			email:    email,
			errMsg:   "user name must be between 3 and 20 characters",
		},
		{
			caseName: "Negative: 無効なメールアドレスの場合はエラーが返る",
			name:     name,
			email:    "invalid",
			errMsg:   "the email format is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			u, err := userDomain.New(tt.name, tt.email)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Nil(t, u)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, u.ID())
				assert.Equal(t, tt.name, u.Name())
				assert.Equal(t, tt.email, u.Email())
			}
		})
	}
}

func TestReconstruct(t *testing.T) {
	var (
		id    = ulid.GenerateStaticULID("user")
		name  = "sato taro"
		email = "sato@example.com"
	)

	t.Run("Positive: ユーザーを再構築できる", func(t *testing.T) {
		u, err := userDomain.Reconstruct(id, name, email)
		assert.NoError(t, err)
		assert.Equal(t, id, u.IDString())
		assert.Equal(t, name, u.Name())
		assert.Equal(t, email, u.Email())
	})
}

func TestChangeName(t *testing.T) {
	var (
		name  = "sato taro"
		email = "sato@example.com"
	)

	tests := []struct {
		caseName string
		newName  string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効な名前の場合は、名前が変更できる",
			newName:  "yamada hanako",
			errMsg:   "",
		},
		{
			caseName: "Negative: 無効な名前の場合はエラーが返る",
			newName:  strings.Repeat("a", 21),
			errMsg:   "user name must be between 3 and 20 characters",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			u, _ := userDomain.New(name, email)
			err := u.ChangeName(tt.newName)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Equal(t, name, u.Name())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newName, u.Name())
			}
		})
	}
}

func TestChangeEmail(t *testing.T) {
	var (
		name  = "sato taro"
		email = "sato@example.com"
	)

	tests := []struct {
		caseName string
		newEmail string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効なメールアドレスの場合は、メールアドレスが変更できる",
			newEmail: "yamada@example.com",
			errMsg:   "",
		},
		{
			caseName: "Negative: 無効なメールアドレスの場合はエラーが返る",
			newEmail: "invalid-email",
			errMsg:   "the email format is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			u, _ := userDomain.New(name, email)
			err := u.ChangeEmail(tt.newEmail)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Equal(t, email, u.Email())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.newEmail, u.Email())
			}
		})
	}
}
