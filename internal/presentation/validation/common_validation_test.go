package validation_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	idVO "github.com/u104rak1/pocgo/internal/domain/value_object/id"
	"github.com/u104rak1/pocgo/internal/presentation/validation"
)

func TestInvalidULID(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効なULID",
			input:    "01FZ7TYX2X5JYZR9DYZXQSTB4H",
			errMsg:   "",
		},
		{
			caseName: "Negative: 空文字列は無効",
			input:    "",
			errMsg:   idVO.ErrInvalidULID.Error(),
		},
		{
			caseName: "Negative: 無効なULID形式",
			input:    "invalid-ulid",
			errMsg:   idVO.ErrInvalidULID.Error(),
		},
		{
			caseName: "Negative: 25文字のULIDは無効",
			input:    "01FZ7TYX2X5JYZR9DYZXQSTB4",
			errMsg:   idVO.ErrInvalidULID.Error(),
		},
		{
			caseName: "Negative: 27文字のULIDは無効",
			input:    "01FZ7TYX2X5JYZR9DYZXQSTB4H1",
			errMsg:   idVO.ErrInvalidULID.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidULID(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidYYYYMMDD(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効な日付",
			input:    "20240301",
			errMsg:   "",
		},
		{
			caseName: "Negative: 非数値の日付は無効",
			input:    "2024030a",
			errMsg:   "must be in a valid format",
		},
		{
			caseName: "Negative: 存在しない日付は無効",
			input:    "20240231",
			errMsg:   "parsing time \"20240231\": day out of range",
		},
		{
			caseName: "Negative: 8桁未満の日付は無効",
			input:    "2024030",
			errMsg:   "must be in a valid format",
		},
		{
			caseName: "Negative: 8桁より多い日付は無効",
			input:    "202403010",
			errMsg:   "must be in a valid format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidYYYYMMDD(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidateDateRange(t *testing.T) {
	tests := []struct {
		caseName string
		from     string
		to       string
		errMsg   string
	}{
		{
			caseName: "Positive: 有効な日付範囲",
			from:     "20240301",
			to:       "20240302",
			errMsg:   "",
		},
		{
			caseName: "Positive: 同じ日付は有効",
			from:     "20240301",
			to:       "20240301",
			errMsg:   "",
		},
		{
			caseName: "Negative: 終了日が開始日より前の日付は無効",
			from:     "20240302",
			to:       "20240301",
			errMsg:   "to date cannot be before from date",
		},
		{
			caseName: "Negative: 開始日が無効な日付形式の場合は無効",
			from:     "2024030",
			to:       "20240301",
			errMsg:   "parsing time \"2024030\" as \"20060102\": cannot parse \"0\" as \"02\"",
		},
		{
			caseName: "Negative: 終了日が無効な日付形式の場合は無効",
			from:     "20240301",
			to:       "2024030",
			errMsg:   "parsing time \"2024030\" as \"20060102\": cannot parse \"0\" as \"02\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidateDateRange(tt.from, tt.to)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidPage(t *testing.T) {
	tests := []struct {
		caseName string
		input    int
		errMsg   string
	}{
		{
			caseName: "Positive: 有効なページ番号",
			input:    1,
			errMsg:   "",
		},
		{
			caseName: "Negative: 0は無効",
			input:    0,
			errMsg:   "page must be greater than 0",
		},
		{
			caseName: "Negative: 負のページ番号は無効",
			input:    -1,
			errMsg:   "page must be greater than 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidPage(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}

func TestValidSort(t *testing.T) {
	tests := []struct {
		caseName string
		input    string
		errMsg   string
	}{
		{
			caseName: "Positive: ASCは有効",
			input:    "ASC",
			errMsg:   "",
		},
		{
			caseName: "Positive: DESCは有効",
			input:    "DESC",
			errMsg:   "",
		},
		{
			caseName: "Negative: 小文字のASCは無効",
			input:    "asc",
			errMsg:   "must be a valid value",
		},
		{
			caseName: "Negative: 小文字のDESCは無効",
			input:    "desc",
			errMsg:   "must be a valid value",
		},
		{
			caseName: "Negative: 無効なソート順は無効",
			input:    "INVALID",
			errMsg:   "must be a valid value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()
			err := validation.ValidSort(tt.input)
			if tt.errMsg == "" {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			}
		})
	}
}
