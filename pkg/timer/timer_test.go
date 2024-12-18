package timer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/pkg/timer"
)

func TestNow(t *testing.T) {
	t.Run("現在時刻をUTCで秒単位に切り捨てて返す", func(t *testing.T) {
		now := timer.Now()
		assert.Equal(t, time.UTC, now.Location())
		assert.Zero(t, now.Nanosecond())
		assert.WithinDuration(t, time.Now().UTC().Truncate(time.Second), now, time.Second)
	})
}

func TestFormatToISO8601(t *testing.T) {
	t.Run("指定された日時をISO8601形式の文字列に変換する", func(t *testing.T) {
		date := time.Date(2022, 3, 15, 10, 30, 0, 0, time.UTC)
		expected := "2022-03-15T10:30:00Z"
		formattedDate := timer.FormatToISO8601(date)
		assert.Equal(t, expected, formattedDate)
	})
}

func TestGetFixedDate(t *testing.T) {
	t.Run("固定の日時を返す", func(t *testing.T) {
		expected := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		fixedDate := timer.GetFixedDate()
		assert.Equal(t, expected, fixedDate)
	})
}

func TestGetFixedDateString(t *testing.T) {
	t.Run("固定の日時をISO8601形式の文字列で返す", func(t *testing.T) {
		expected := "2021-01-01T00:00:00Z"
		fixedDateString := timer.GetFixedDateString()
		assert.Equal(t, expected, fixedDateString)
	})
}

func TestParseYYYYMMDD(t *testing.T) {
	tests := []struct {
		caseName string
		dateStr  string
		want     time.Time
		errMsg   string
	}{
		{
			caseName: "YYYYMMDD形式の日付文字列をパースできる",
			dateStr:  "20240101",
			want:     time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			errMsg:   "",
		},
		{
			caseName: "不正な形式の日付文字列はパースに失敗する",
			dateStr:  "2024-01-01",
			want:     time.Time{},
			errMsg:   "invalid date format: 2024-01-01, expected YYYYMMDD",
		},
		{
			caseName: "日付として解釈できない文字列はパースに失敗する",
			dateStr:  "invalid-date",
			want:     time.Time{},
			errMsg:   "invalid date format: invalid-date, expected YYYYMMDD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			got, err := timer.ParseYYYYMMDD(tt.dateStr)

			if tt.errMsg != "" {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
				assert.Equal(t, tt.want, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
