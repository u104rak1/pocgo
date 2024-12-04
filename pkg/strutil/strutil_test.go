package strutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104raki/pocgo/pkg/strutil"
)

func TestToSnakeFromCamel(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		expected string
	}{
		{
			caseName: "Converts CamelCase to snake_case.",
			input:    "CamelCase",
			expected: "camel_case",
		},
		{
			caseName: "Handles single word.",
			input:    "Word",
			expected: "word",
		},
		{
			caseName: "Handles empty string.",
			input:    "",
			expected: "",
		},
		{
			caseName: "Handles already snake_case.",
			input:    "snake_case",
			expected: "snake_case",
		},
		{
			caseName: "Handles mixed case.",
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
			caseName: "Converts spaces to kebab-case.",
			input:    "hello world",
			expected: "hello-world",
		},
		{
			caseName: "Handles multiple spaces.",
			input:    "go is awesome",
			expected: "go-is-awesome",
		},
		{
			caseName: "Handles leading and trailing spaces.",
			input:    "  leading and trailing  ",
			expected: "leading-and-trailing",
		},
		{
			caseName: "Handles empty string.",
			input:    "",
			expected: "",
		},
		{
			caseName: "Handles no spaces.",
			input:    "kebabcase",
			expected: "kebabcase",
		},
		{
			caseName: "Handles special characters.",
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
