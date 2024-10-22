package timer

import "time"

func Now() time.Time {
	return time.Now().UTC().Truncate(time.Second)
}
