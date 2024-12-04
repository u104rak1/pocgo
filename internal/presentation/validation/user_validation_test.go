package validation_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	authDomain "github.com/u104raki/pocgo/internal/domain/authentication"
	userDomain "github.com/u104raki/pocgo/internal/domain/user"
	"github.com/u104raki/pocgo/internal/presentation/validation"
)

func TestValidUserName(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		wantErr  string
	}{
		{
			caseName: "An empty name is invalid.",
			input:    "",
			wantErr:  "cannot be blank",
		},
		{
			caseName: "A name less than 2-characters are invalid.",
			input:    strings.Repeat("a", userDomain.NameMinLength-1),
			wantErr:  "the length must be between 3 and 20",
		},
		{
			caseName: "Valid name must be a minimum 3-characters.",
			input:    strings.Repeat("a", userDomain.NameMinLength),
			wantErr:  "",
		},
		{
			caseName: "Valid name must be a maximum of 20-characters.",
			input:    strings.Repeat("a", userDomain.NameMaxLength),
			wantErr:  "",
		},
		{
			caseName: "A name longer than 21-characters are invalid.",
			input:    strings.Repeat("a", userDomain.NameMaxLength+1),
			wantErr:  "the length must be between 3 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserName(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidUserEmail(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		wantErr  string
	}{
		{
			caseName: "An empty email is invalid.",
			input:    "",
			wantErr:  "cannot be blank",
		},
		{
			caseName: "An invalid email format is invalid.",
			input:    "invalid",
			wantErr:  "the email format is invalid",
		},
		{
			caseName: "A valid email format is valid.",
			input:    "test@example.com",
			wantErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserEmail(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}

func TestValidUserPassword(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		wantErr  string
	}{
		{
			caseName: "An empty password is invalid.",
			input:    "",
			wantErr:  "cannot be blank",
		},
		{
			caseName: "A password less than 7-characters are invalid.",
			input:    strings.Repeat("a", authDomain.PasswordMinLength-1),
			wantErr:  "the length must be between 8 and 20",
		},
		{
			caseName: "Valid password must be a minimum 8-characters.",
			input:    strings.Repeat("a", authDomain.PasswordMinLength),
			wantErr:  "",
		},
		{
			caseName: "Valid password must be a maximum of 20-characters.",
			input:    strings.Repeat("a", authDomain.PasswordMaxLength),
			wantErr:  "",
		},
		{
			caseName: "A password longer than 21-characters are invalid.",
			input:    strings.Repeat("a", authDomain.PasswordMaxLength+1),
			wantErr:  "the length must be between 8 and 20",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidUserPassword(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}
