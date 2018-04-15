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
		transaction string
		exp         string
	}{
		// TODO
	}
	for _, tc := range tcs {
		got, err := NewRow(tc.transaction, account)
		if err != nil {
			t.Errorf("%v", err)
		}
		if got.String() != tc.exp {
			t.Errorf("expect csv %q, got %v", tc.exp, got)
		}
	}
}
