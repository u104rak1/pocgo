package strutil

import (
	"bytes"
	"strings"
)

// キャメルケースをスネークケースに変換します。
func ToSnakeFromCamel(s string) string {
	var result bytes.Buffer
	for i, char := range s {
		if i > 0 && char >= 'A' && char <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(char)
	}
	return strings.ToLower(result.String())
}

// スペースを含む文字列をケバブ文字列に変換します。前後のスペースはトリミングされます。
func ToKebabFromSpace(s string) string {
	trimmed := strings.TrimSpace(s)
	return strings.ReplaceAll(trimmed, " ", "-")
}

// 文字列のポインタを返します。
func StrPointer(s string) *string {
	return &s
}
