package date

import (
	"os"
	"testing"
	"time"
)

const longForm = "Mon Jan 2 15:04:05 -0700 MST 2006"

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestDate(t *testing.T) {
	tcs := []struct {
		second uint
		expect string
	}{
		{
			1,
			"Mon Jan 1 00:00:01 +0000 UTC 2000",
		},
	}
	for _, tc := range tcs {
		got := New(tc.second)
		exp, _ := time.Parse(longForm, tc.expect)
		if !got.Equal(exp) {
			t.Fatalf("expect ripple epoch %d equal %q, got %v", tc.second, tc.expect, got)
		}
	}
}
