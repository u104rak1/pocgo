package id_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
)

func TestNewUserID(t *testing.T) {
	t.Run("新規UserIDが生成され、有効なULIDフォーマットであること", func(t *testing.T) {
		id := idVO.NewUserID()
		assert.NotEmpty(t, id.String())
		assert.True(t, id.IsValid())
	})
}

func TestUserIDFromString(t *testing.T) {
	t.Run("有効なULIDからUserIDを生成できること", func(t *testing.T) {
		input := "01H2X5JMIN3P8T68PYHXXVK5XN"
		id, err := idVO.UserIDFromString(input)
		assert.NoError(t, err)
		assert.Equal(t, input, id.String())
	})
}

func TestNewUserIDForTest(t *testing.T) {
	tests := []struct {
		name     string
		seed1    string
		seed2    string
		wantSame bool
	}{
		{
			name:     "同じシードから同じUserIDが生成されること",
			seed1:    "test-user-1",
			seed2:    "test-user-1",
			wantSame: true,
		},
		{
			name:     "異なるシードから異なるUserIDが生成されること",
			seed1:    "test-user-1",
			seed2:    "test-user-2",
			wantSame: false,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			id1 := idVO.NewUserIDForTest(tt.seed1)
			id2 := idVO.NewUserIDForTest(tt.seed2)

			assert.Equal(t, tt.wantSame, id1.Equals(id2))
			assert.True(t, id1.IsValid())
			assert.True(t, id2.IsValid())
		})
	}
}
