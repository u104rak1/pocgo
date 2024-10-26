package strutil

import (
	"bytes"
	"strings"
)

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
