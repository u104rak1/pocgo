package timer

import (
	"fmt"
	"time"
)

// Now returns the current UTC time truncated to seconds
func Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

// FormatToISO8601 converts a time.Time value to an ISO8601 string in UTC.
func FormatToISO8601(t time.Time) string {
	return t.UTC().Truncate(time.Second).Format(time.RFC3339)
}

// Default fixed date: "2021-01-01T00:00:00Z".
var fixedDate = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

// GetFixedDate returns the current fixed date for use in testing.
// Default value is "2021-01-01T00:00:00Z".
func GetFixedDate() time.Time {
	return fixedDate
}

// GetFixedDateString returns the fixed date as an ISO8601 formatted string.
// Default value is "2021-01-01T00:00:00Z".
func GetFixedDateString() string {
	return fixedDate.Format(time.RFC3339)
}

// ParseYYYYMMDD parses a date string in the format "YYYYMMDD" into a time.Time value in UTC.
// The result is truncated to seconds.
func ParseYYYYMMDD(dateStr string) (time.Time, error) {
	parsedTime, err := time.Parse("20060102", dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format: %s, expected YYYYMMDD", dateStr)
	}
	return parsedTime.UTC().Truncate(time.Second), nil
}
