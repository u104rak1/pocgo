package id_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

func TestNewAccountID(t *testing.T) {
	t.Run("新規AccountIDが生成され、有効なULIDフォーマットであること", func(t *testing.T) {
		id := idVO.NewAccountID()
		assert.NotEmpty(t, id.String())
		assert.True(t, id.IsValid())
	})
}

func TestAccountIDFromString(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		errMsg string
	}{
		{
			name:   "Positive: 有効なULIDからAccountIDを生成できること",
			input:  "01H2X5JMIN3P8T68PYHXXVK5XN",
			errMsg: "",
		},
		{
			name:   "Negative: 不正なULIDからAccountIDを生成できないこと",
			input:  "invalid-ulid",
			errMsg: "invalid account id: invalid ulid",
		},
		{
			name:   "Negative: 空文字列からAccountIDを生成できないこと",
			input:  "",
			errMsg: "invalid account id: id must not be empty",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id, err := idVO.AccountIDFromString(tt.input)
			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				return
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.input, id.String())
			}
		})
	}
}

func TestNewAccountIDForTest(t *testing.T) {
	tests := []struct {
		name     string
		seed1    string
		seed2    string
		wantSame bool
	}{
		{
			name:     "同じシードから同じAccountIDが生成されること",
			seed1:    "test-account-1",
			seed2:    "test-account-1",
			wantSame: true,
		},
		{
			name:     "異なるシードから異なるAccountIDが生成されること",
			seed1:    "test-account-1",
			seed2:    "test-account-2",
			wantSame: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id1 := idVO.NewAccountIDForTest(tt.seed1)
			id2 := idVO.NewAccountIDForTest(tt.seed2)

			assert.Equal(t, tt.wantSame, id1.Equals(id2))
			assert.True(t, id1.IsValid())
			assert.True(t, id2.IsValid())
		})
	}
}
