package timer_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ucho456job/pocgo/pkg/timer"
)

func TestNow(t *testing.T) {
	t.Run("Successfully returns current time in UTC truncated to seconds.", func(t *testing.T) {
		now := timer.Now()
		assert.Equal(t, time.UTC, now.Location())
		assert.Zero(t, now.Nanosecond())
		assert.WithinDuration(t, time.Now().UTC().Truncate(time.Second), now, time.Second)
	})
}
