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
			caseName: "Valid: シンプルなメールアドレス形式",
			email:    "test@example.com",
			want:     true,
		},
		{
			caseName: "Valid: サブドメインを含むメールアドレス",
			email:    "user@mail.example.co.jp",
			want:     true,
		},
		{
			caseName: "Valid: プラス記号を含むメールアドレス",
			email:    "user+123@example.com",
			want:     true,
		},
		{
			caseName: "Valid: ハイフンを含むドメインのメールアドレス",
			email:    "user@my-example.com",
			want:     true,
		},
		{
			caseName: "Invalid: @シンボルがない",
			email:    "userexample.com",
			want:     false,
		},
		{
			caseName: "Invalid: @シンボルが2つある",
			email:    "user@@example.com",
			want:     false,
		},
		{
			caseName: "Invalid: ドメインが欠落している",
			email:    "user@.com",
			want:     false,
		},
		{
			caseName: "Invalid: 不正な文字を含むドメイン",
			email:    "user@exa!mple.com",
			want:     false,
		},
		{
			caseName: "Invalid: ドメインが短すぎる",
			email:    "user@example.c",
			want:     false,
		},
		{
			caseName: "Invalid: トップレベルドメインがない",
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
