package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/presentation/shared/validation"
	ulidUtil "github.com/ucho456job/pocgo/pkg/ulid"
)

func TestInvalidULID(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			"Valid ULID",
			"01FZ7TYX2X5JYZR9DYZXQSTB4H",
			"",
		},
		{
			"Empty string is invalid",
			"",
			ulidUtil.ErrInvalidULID.Error(),
		},
		{
			"Invalid ULID format",
			"invalid-ulid",
			ulidUtil.ErrInvalidULID.Error(),
		},
		{
			"Too short ULID (25 characters)",
			"01FZ7TYX2X5JYZR9DYZXQSTB4",
			ulidUtil.ErrInvalidULID.Error(),
		},
		{
			"Too long ULID (27 characters)",
			"01FZ7TYX2X5JYZR9DYZXQSTB4H1",
			ulidUtil.ErrInvalidULID.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidULID(tt.input)
			if tt.wantErr == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err.Error())
			}
		})
	}
}
