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
			caseName: "基本的なメールアドレスを検証できる",
			email:    "test@example.com",
			want:     true,
		},
		{
			caseName: "サブドメインを含むメールアドレスを検証できる",
			email:    "user@mail.example.co.jp",
			want:     true,
		},
		{
			caseName: "プラス記号を含むメールアドレスを検証できる",
			email:    "user+123@example.com",
			want:     true,
		},
		{
			caseName: "ハイフンを含むメールアドレスを検証できる",
			email:    "user@my-example.com",
			want:     true,
		},
		{
			caseName: "アットマーク(@)がないメールアドレスは無効",
			email:    "userexample.com",
			want:     false,
		},
		{
			caseName: "アットマーク(@)が2つあるメールアドレスは無効",
			email:    "user@@example.com",
			want:     false,
		},
		{
			caseName: "ドメインがないメールアドレスは無効",
			email:    "user@.com",
			want:     false,
		},
		{
			caseName: "無効な文字が含まれるメールアドレスは無効",
			email:    "user@exa!mple.com",
			want:     false,
		},
		{
			caseName: "ドメインが短すぎるメールアドレスは無効",
			email:    "user@example.c",
			want:     false,
		},
		{
			caseName: "トップレベルドメインがないメールアドレスは無効",
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
