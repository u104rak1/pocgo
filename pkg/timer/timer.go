package timer

import "time"

func Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}

func FormatToISO8601(t time.Time) string {
	return t.UTC().Truncate(time.Second).Format(time.RFC3339)
}

var fixedDate = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

// Return a fixed date for use in testing.
// Return time.Time of "2021-01-01T00:00:00Z".
func GetFixedDate() time.Time {
	return fixedDate
}

// Returns a fixed date for use in testing.
// Returns a string in ISO8601 format of "2021-01-01T00:00:00Z".
func GetFixedDateString() string {
	return fixedDate.Format(time.RFC3339)
}
