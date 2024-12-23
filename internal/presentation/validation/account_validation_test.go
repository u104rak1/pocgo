package validation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	accountDomain "github.com/u104rak1/pocgo/internal/domain/account"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
)

func TestValidAccountName(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効な口座名",
			input:    "test",
			errMsg:   "",
		},
		{
			caseName: "Negative: 空文字列は無効",
			input:    "",
			errMsg:   "cannot be blank",
		},
		{
			caseName: "Negative: 2文字未満は無効",
			input:    strings.Repeat("a", accountDomain.NameMinLength-1),
			errMsg:   "the length must be between 3 and 20",
		},
		{
			caseName: "Positive: 3文字以上は有効",
			input:    strings.Repeat("a", accountDomain.NameMinLength),
			errMsg:   "",
		},
		{
			caseName: "Positive: 20文字以下は有効",
			input:    strings.Repeat("a", accountDomain.NameMaxLength),
			errMsg:   "",
		},
		{
			caseName: "Negative: 21文字以上は無効",
			input:    strings.Repeat("a", accountDomain.NameMaxLength+1),
			errMsg:   "the length must be between 3 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidAccountName(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidAccountPassword(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Negative: 空文字列は無効",
			input:    "",
			errMsg:   "cannot be blank",
		},
		{
			caseName: "Negative: 3文字未満は無効",
			input:    strings.Repeat("a", accountDomain.PasswordLength-1),
			errMsg:   "the length must be exactly 4",
		},
		{
			caseName: "Positive: 4文字は有効なパスワード",
			input:    strings.Repeat("a", accountDomain.PasswordLength),
			errMsg:   "",
		},
		{
			caseName: "Negative: 5文字以上は無効",
			input:    strings.Repeat("a", accountDomain.PasswordLength+1),
			errMsg:   "the length must be exactly 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidAccountPassword(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}
