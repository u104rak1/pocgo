package strutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/pkg/strutil"
)

func TestToSnakeFromCamel(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		expected string
	}{
		{
			caseName: "キャメルケースをスネークケースに変換できる",
			input:    "CamelCase",
			expected: "camel_case",
		},
		{
			caseName: "単一の単語を処理できる",
			input:    "Word",
			expected: "word",
		},
		{
			caseName: "空文字を処理できる",
			input:    "",
			expected: "",
		},
		{
			caseName: "すでにスネークケースの文字列を処理できる",
			input:    "snake_case",
			expected: "snake_case",
		},
		{
			caseName: "大文字小文字が混在する文字列を処理できる",
			input:    "MixedCASEExample",
			expected: "mixed_c_a_s_e_example",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			result := strutil.ToSnakeFromCamel(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToKebabFromSpace(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		expected string
	}{
		{
			caseName: "スペースをケバブケースに変換できる",
			input:    "hello world",
			expected: "hello-world",
		},
		{
			caseName: "複数のスペースを処理できる",
			input:    "go is awesome",
			expected: "go-is-awesome",
		},
		{
			caseName: "前後のスペースを処理できる",
			input:    "  leading and trailing  ",
			expected: "leading-and-trailing",
		},
		{
			caseName: "空文字を処理できる",
			input:    "",
			expected: "",
		},
		{
			caseName: "スペースがない文字列を処理できる",
			input:    "kebabcase",
			expected: "kebabcase",
		},
		{
			caseName: "特殊文字を含む文字列を処理できる",
			input:    "hello @world!",
			expected: "hello-@world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			result := strutil.ToKebabFromSpace(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestStrPointer(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
	}{
		{
			caseName: "通常の文字列のポインタを取得できる",
			input:    "hello",
		},
		{
			caseName: "空文字列のポインタを取得できる",
			input:    "",
		},
		{
			caseName: "日本語文字列のポインタを取得できる",
			input:    "こんにちは",
		},
		{
			caseName: "特殊文字を含む文字列のポインタを取得できる",
			input:    "hello@world!",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			result := strutil.StrPointer(tt.input)
			assert.NotNil(t, result)
			assert.Equal(t, tt.input, *result)
		})
	}
}
