package timer

import (
	"fmt"
	"time"
)

// Now は現在のUTC時刻を秒単位で切り捨てて返します
func Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

// FormatToISO8601 は time.Time 値をUTCのISO8601形式の文字列に変換します
func FormatToISO8601(t time.Time) string {
	return t.UTC().Truncate(time.Second).Format(time.RFC3339)
}

var fixedDate = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

// GetFixedDate はテストで使用する固定の日時を返します
// デフォルト値は "2021-01-01T00:00:00Z" です
func GetFixedDate() time.Time {
	return fixedDate
}

// GetFixedDateString は固定の日時をISO8601形式の文字列で返します
// デフォルト値は "2021-01-01T00:00:00Z" です
func GetFixedDateString() string {
	return fixedDate.Format(time.RFC3339)
}

// ParseYYYYMMDD は "YYYYMMDD" 形式の日付文字列をUTCのtime.Time値にパースします
// 結果は秒単位で切り捨てられます
func ParseYYYYMMDD(dateStr string) (time.Time, error) {
	parsedTime, err := time.Parse("20060102", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %s, expected YYYYMMDD", dateStr)
	}
	return parsedTime.UTC().Truncate(time.Second), nil
}
