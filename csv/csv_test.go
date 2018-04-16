package csv

import (
	"os"
	"testing"
)

const (
	account = ""
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestToString(t *testing.T) {
	tcs := []struct {
		format      string
		transaction string
		exp         string
	}{
		// TODO
	}
	for _, tc := range tcs {
		got := Factory[tc.format]()
		if err := got.New(tc.transaction, account); err != nil {
			t.Errorf("%v", err)
		}
		if got.String() != tc.exp {
			t.Errorf("expect csv %q, got %v", tc.exp, got)
		}
	}
}
