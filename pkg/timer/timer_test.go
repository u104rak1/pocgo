package timer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/u104rak1/pocgo/pkg/timer"
)

func TestNow(t *testing.T) {
	t.Run("Successfully returns current time in UTC, truncated to the nearest second.", func(t *testing.T) {
		now := timer.Now()
		assert.Equal(t, time.UTC, now.Location())
		assert.Zero(t, now.Nanosecond())
		assert.WithinDuration(t, time.Now().UTC().Truncate(time.Second), now, time.Second)
	})
}

func TestFormatToISO8601(t *testing.T) {
	t.Run("Convert the specified date and time to a string in iso8601 format.", func(t *testing.T) {
		date := time.Date(2022, 3, 15, 10, 30, 0, 0, time.UTC)
		expected := "2022-03-15T10:30:00Z"
		formattedDate := timer.FormatToISO8601(date)
		assert.Equal(t, expected, formattedDate)
	})
}

func TestGetFixedDate(t *testing.T) {
	t.Run("Returns a fixed date.", func(t *testing.T) {
		expected := time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
		fixedDate := timer.GetFixedDate()
		assert.Equal(t, expected, fixedDate)
	})
}

func TestGetFixedDateString(t *testing.T) {
	t.Run("Returns a string in ISO8601 format with a fixed date.", func(t *testing.T) {
		expected := "2021-01-01T00:00:00Z"
		fixedDateString := timer.GetFixedDateString()
		assert.Equal(t, expected, fixedDateString)
	})
}

func TestParseYYYYMMDD(t *testing.T) {
	t.Run("Successfully parses a valid YYYYMMDD date string.", func(t *testing.T) {
		dateStr := "20240101"
		expected := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		parsedDate, err := timer.ParseYYYYMMDD(dateStr)
		assert.NoError(t, err)
		assert.Equal(t, expected, parsedDate)
	})

	t.Run("Fails to parse an invalid YYYYMMDD date string.", func(t *testing.T) {
		invalidDateStr := "2024-01-01"
		_, err := timer.ParseYYYYMMDD(invalidDateStr)
		assert.Error(t, err)
	})

	t.Run("Fails to parse a malformed date string.", func(t *testing.T) {
		malformedDateStr := "invalid-date"
		_, err := timer.ParseYYYYMMDD(malformedDateStr)
		assert.Error(t, err)
	})
}
