package validation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	authDomain "github.com/u104rak1/pocgo/internal/domain/authentication"
	userDomain "github.com/u104rak1/pocgo/internal/domain/user"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
)

func TestValidUserName(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効なユーザー名",
			input:    "test",
			errMsg:   "",
		},
		{
			caseName: "Negative: 空の名前は無効",
			input:    "",
			errMsg:   "cannot be blank",
		},
		{
			caseName: "Negative: 2文字未満の名前は無効",
			input:    strings.Repeat("a", userDomain.NameMinLength-1),
			errMsg:   "the length must be between 3 and 20",
		},
		{
			caseName: "Positive: 3文字以上の名前は有効",
			input:    strings.Repeat("a", userDomain.NameMinLength),
			errMsg:   "",
		},
		{
			caseName: "Positive: 20文字以下の名前は有効",
			input:    strings.Repeat("a", userDomain.NameMaxLength),
			errMsg:   "",
		},
		{
			caseName: "Negative: 21文字以上の名前は無効",
			input:    strings.Repeat("a", userDomain.NameMaxLength+1),
			errMsg:   "the length must be between 3 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserName(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidUserEmail(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効なメールアドレスは有効",
			input:    "test@example.com",
			errMsg:   "",
		},
		{
			caseName: "Negative: 空のメールアドレスは無効",
			input:    "",
			errMsg:   "cannot be blank",
		},
		{
			caseName: "Negative: 無効なメールアドレスは無効",
			input:    "invalid",
			errMsg:   "the email format is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserEmail(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidUserPassword(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効なパスワード",
			input:    "password",
			errMsg:   "",
		},
		{
			caseName: "Negative: 空のパスワードは無効",
			input:    "",
			errMsg:   "cannot be blank",
		},
		{
			caseName: "Negative: 7文字未満のパスワードは無効",
			input:    strings.Repeat("a", authDomain.PasswordMinLength-1),
			errMsg:   "the length must be between 8 and 20",
		},
		{
			caseName: "Positive: 8文字以上のパスワードは有効",
			input:    strings.Repeat("a", authDomain.PasswordMinLength),
			errMsg:   "",
		},
		{
			caseName: "Positive: 20文字以下のパスワードは有効",
			input:    strings.Repeat("a", authDomain.PasswordMaxLength),
			errMsg:   "",
		},
		{
			caseName: "Negative: 21文字以上のパスワードは無効",
			input:    strings.Repeat("a", authDomain.PasswordMaxLength+1),
			errMsg:   "the length must be between 8 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserPassword(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}
