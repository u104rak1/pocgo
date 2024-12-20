package id_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

// テスト用のマーカー型
type testIDType struct{}

func TestNew(t *testing.T) {
	t.Run("新規IDが生成され、有効なULIDフォーマットであること", func(t *testing.T) {
		id := idVO.New[testIDType]()
		assert.NotEmpty(t, id.String())
		assert.True(t, id.IsValid())
	})
}

func TestNewFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr error
	}{
		{
			name:    "Positive: 有効なULID",
			input:   "01H2X5JMIN3P8T68PYHXXVK5XN",
			wantErr: nil,
		},
		{
			name:    "Negative: 空文字列",
			input:   "",
			wantErr: idVO.ErrEmptyID,
		},
		{
			name:    "Negative: 無効なULID形式",
			input:   "invalid-ulid",
			wantErr: idVO.ErrInvalidULID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := idVO.NewFromString[testIDType](tt.input)
			if tt.wantErr != nil {
				assert.Equal(t, tt.wantErr, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.input, id.String())
		})
	}
}

func TestID_Equals(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		id1      string
		id2      string
		expected bool
	}{
		{
			name:     "同じID",
			id1:      "01H2X5JMIN3P8T68PYHXXVK5XN",
			id2:      "01H2X5JMIN3P8T68PYHXXVK5XN",
			expected: true,
		},
		{
			name:     "異なるID",
			id1:      "01H2X5JMIN3P8T68PYHXXVK5XN",
			id2:      "09H2X5JMIN3P8T68PYHXXVK5XM",
			expected: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id1, err := idVO.NewFromString[testIDType](tt.id1)
			assert.NoError(t, err)
			id2, err := idVO.NewFromString[testIDType](tt.id2)
			assert.NoError(t, err)

			assert.Equal(t, tt.expected, id1.Equals(id2))
		})
	}
}

// IsValidのテストケースはNewFromStringのテストケースで確認済みなので不個別のテストは不要
// func TestID_IsValid(t *testing.T) {}

func TestNewForTest(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		seed1    string
		seed2    string
		wantSame bool
	}{
		{
			name:     "同じシードから同じIDが生成されること",
			seed1:    "test-seed",
			seed2:    "test-seed",
			wantSame: true,
		},
		{
			name:     "異なるシードから異なるIDが生成されること",
			seed1:    "test-seed-1",
			seed2:    "test-seed-2",
			wantSame: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id1 := idVO.NewForTest[testIDType](tt.seed1)
			id2 := idVO.NewForTest[testIDType](tt.seed2)

			assert.Equal(t, tt.wantSame, id1.Equals(id2))
			assert.True(t, id1.IsValid())
			assert.True(t, id2.IsValid())
		})
	}
}
