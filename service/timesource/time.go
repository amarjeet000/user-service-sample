package timesource

import "time"

func CurrentTime() time.Time {
	return time.Now().UTC()
}
