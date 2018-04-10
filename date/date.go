package date

import "time"

var (
	epoch = time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC)
)

func New(sec uint32) time.Time {
	return epoch.Add(time.Second * time.Duration(sec))
}
