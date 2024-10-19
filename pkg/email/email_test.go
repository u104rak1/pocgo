package email_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/pkg/email"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		caseName string
		email    string
		want     bool
	}{
		{
			caseName: "Valid: return true, if the email is valid.",
			email:    "test@example.com",
			want:     true,
		},
		{
			caseName: "Valid: return true, if a valid email address including any subdomains.",
			email:    "user@mail.example.co.jp",
			want:     true,
		},
		{
			caseName: "Valid: return true, if the email address contains a plus sign.",
			email:    "user+123@example.com",
			want:     true,
		},
		{
			caseName: "Valid: return true, if the email address contains a hyphen.",
			email:    "user@my-example.com",
			want:     true,
		},
		{
			caseName: "Invalid: return false, if there is no @.",
			email:    "userexample.com",
			want:     false,
		},
		{
			caseName: "Invalid: return false, if there are two @ signs.",
			email:    "user@@example.com",
			want:     false,
		},
		{
			caseName: "Invalid: return false, if there is no domain.",
			email:    "user@.com",
			want:     false,
		},
		{
			caseName: "Invalid: return false, if there is a invalid character.",
			email:    "user@exa!mple.com",
			want:     false,
		},
		{
			caseName: "Invalid: return false, if domain is too short.",
			email:    "user@example.c",
			want:     false,
		},
		{
			caseName: "Invalid: return false, if there is no top-level domain.",
			email:    "user@example",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			got := email.IsValid(tt.email)
			assert.Equal(t, tt.want, got)
		})
	}
}
