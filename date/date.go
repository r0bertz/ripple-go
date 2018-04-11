package date

import "time"

var (
	// Epoch is ripple epoch.
	Epoch = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
)

// New returns a time.Time object from a number which is the number of seconds
// since the "Ripple Epoch" of January 1, 2000 (00:00 UTC).
func New(sec uint32) time.Time {
	return Epoch.Add(time.Second * time.Duration(sec))
}
