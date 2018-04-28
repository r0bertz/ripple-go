package csv

import (
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestFormat(t *testing.T) {
	tcs := []struct {
		row    Row
		format string
		exp    string
	}{
		// TODO
	}
	for _, tc := range tcs {
		formatter := FormatterFactory[tc.format]()
		got, err := formatter.Format(tc.row)
		if err != nil {
			t.Errorf("%v", err)
		}
		if got != tc.exp {
			t.Errorf("expect csv %q, got %v", tc.exp, got)
		}
	}
}
