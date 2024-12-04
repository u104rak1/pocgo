package ulid_test

import (
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	myUlid "github.com/u104raki/pocgo/pkg/ulid"
)

func TestNew(t *testing.T) {
	t.Run("Successfully generates a new ULID in the correct format.", func(t *testing.T) {
		id := myUlid.New()
		assert.NotEmpty(t, id)
		_, err := ulid.Parse(id)
		assert.NoError(t, err)
	})
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		caseName string
		ulid     string
		want     bool
	}{
		{
			caseName: "Successfully returns true if the ULID is valid.",
			ulid:     myUlid.New(),
			want:     true,
		},
		{
			caseName: "Fails to validate ULID, returns false if the ULID is invalid.",
			ulid:     "invalid-ulid",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			got := myUlid.IsValid(tt.ulid)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenerateStaticULID(t *testing.T) {
	t.Run("Successfully returns the same ULID for identical seeds.", func(t *testing.T) {
		seed := "test-seed"
		ulid1 := myUlid.GenerateStaticULID(seed)
		ulid2 := myUlid.GenerateStaticULID(seed)
		assert.Equal(t, ulid1, ulid2)
	})
	t.Run("Successfully returns different ULIDs for different seeds.", func(t *testing.T) {
		seed1 := "test-seed-1"
		seed2 := "test-seed-2"
		ulid1 := myUlid.GenerateStaticULID(seed1)
		ulid2 := myUlid.GenerateStaticULID(seed2)
		assert.NotEqual(t, ulid1, ulid2)
	})
}
