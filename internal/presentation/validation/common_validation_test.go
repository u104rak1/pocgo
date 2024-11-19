package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/internal/presentation/validation"
	ulidUtil "github.com/ucho456job/pocgo/pkg/ulid"
)

func TestInvalidULID(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		wantErr  string
	}{
		{
			caseName: "Valid ULID",
			input:    "01FZ7TYX2X5JYZR9DYZXQSTB4H",
			wantErr:  "",
		},
		{
			caseName: "Empty string is invalid",
			input:    "",
			wantErr:  ulidUtil.ErrInvalidULID.Error(),
		},
		{
			caseName: "Invalid ULID format",
			input:    "invalid-ulid",
			wantErr:  ulidUtil.ErrInvalidULID.Error(),
		},
		{
			caseName: "Too short ULID (25 characters)",
			input:    "01FZ7TYX2X5JYZR9DYZXQSTB4",
			wantErr:  ulidUtil.ErrInvalidULID.Error(),
		},
		{
			caseName: "Too long ULID (27 characters)",
			input:    "01FZ7TYX2X5JYZR9DYZXQSTB4H1",
			wantErr:  ulidUtil.ErrInvalidULID.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
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
