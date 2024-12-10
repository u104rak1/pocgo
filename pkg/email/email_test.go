package email_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/pkg/email"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		caseName string
		email    string
		want     bool
	}{
		{
			caseName: "Successfully returns true if the email is valid.",
			email:    "test@example.com",
			want:     true,
		},
		{
			caseName: "Successfully returns true if a valid email address includes any subdomains.",
			email:    "user@mail.example.co.jp",
			want:     true,
		},
		{
			caseName: "Successfully returns true if the email address contains a plus sign.",
			email:    "user+123@example.com",
			want:     true,
		},
		{
			caseName: "Successfully returns true if the email address contains a hyphen.",
			email:    "user@my-example.com",
			want:     true,
		},
		{
			caseName: "Fails to validate email, returns false if there is no '@'.",
			email:    "userexample.com",
			want:     false,
		},
		{
			caseName: "Fails to validate email, returns false if there are two '@' signs.",
			email:    "user@@example.com",
			want:     false,
		},
		{
			caseName: "Fails to validate email, returns false if there is no domain.",
			email:    "user@.com",
			want:     false,
		},
		{
			caseName: "Fails to validate email, returns false if there is an invalid character.",
			email:    "user@exa!mple.com",
			want:     false,
		},
		{
			caseName: "Fails to validate email, returns false if the domain is too short.",
			email:    "user@example.c",
			want:     false,
		},
		{
			caseName: "Fails to validate email, returns false if there is no top-level domain.",
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
