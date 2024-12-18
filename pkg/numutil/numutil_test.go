package numutil_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/pkg/numutil"
)

func TestIntPointer(t *testing.T) {
	tests := []struct {
		caseName string
		input    int
	}{
		{
			caseName: "正の整数のポインタを取得できる",
			input:    42,
		},
		{
			caseName: "負の整数のポインタを取得できる",
			input:    -10,
		},
		{
			caseName: "ゼロのポインタを取得できる",
			input:    0,
		},
		{
			caseName: "最大値のポインタを取得できる",
			input:    int(^uint(0) >> 1),
		},
		{
			caseName: "最小値のポインタを取得できる",
			input:    -int(^uint(0)>>1) - 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			result := numutil.IntPointer(tt.input)
			assert.NotNil(t, result)
			assert.Equal(t, tt.input, *result)
		})
	}
}
