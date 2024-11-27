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

func TestValidYYYYMMDD(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		wantErr  bool
	}{
		{
			caseName: "Valid date",
			input:    "20240301",
			wantErr:  false,
		},
		{
			caseName: "Invalid format (non-numeric)",
			input:    "2024030a",
			wantErr:  true,
		},
		{
			caseName: "Invalid date (non-existent date)",
			input:    "20240231",
			wantErr:  true,
		},
		{
			caseName: "Insufficient length",
			input:    "2024030",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidYYYYMMDD(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateDateRange(t *testing.T) {
	tests := []struct {
		caseName string
		from     string
		to       string
		wantErr  bool
	}{
		{
			caseName: "Valid date range",
			from:     "20240301",
			to:       "20240302",
			wantErr:  false,
		},
		{
			caseName: "Same date",
			from:     "20240301",
			to:       "20240301",
			wantErr:  false,
		},
		{
			caseName: "Invalid range (end date before start date)",
			from:     "20240302",
			to:       "20240301",
			wantErr:  true,
		},
		{
			caseName: "Invalid date format (from)",
			from:     "2024030",
			to:       "20240301",
			wantErr:  true,
		},
		{
			caseName: "Invalid date format (to)",
			from:     "20240301",
			to:       "2024030",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidateDateRange(tt.from, tt.to)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidPage(t *testing.T) {
	tests := []struct {
		caseName string
		input    int
		wantErr  bool
	}{
		{
			caseName: "Valid page number",
			input:    1,
			wantErr:  false,
		},
		{
			caseName: "Invalid page number (zero)",
			input:    0,
			wantErr:  true,
		},
		{
			caseName: "Invalid page number (negative)",
			input:    -1,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidPage(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidSort(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		wantErr  bool
	}{
		{
			caseName: "Valid sort order (ASC)",
			input:    "ASC",
			wantErr:  false,
		},
		{
			caseName: "Valid sort order (DESC)",
			input:    "DESC",
			wantErr:  false,
		},
		{
			caseName: "Invalid sort order (lowercase)",
			input:    "asc",
			wantErr:  true,
		},
		{
			caseName: "Invalid sort order (invalid value)",
			input:    "INVALID",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidSort(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
