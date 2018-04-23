package csv

import (
	"log"
	"os"
	"testing"

	"github.com/r0bertz/ripple/data"
)

const (
	account = "rrrrrrrrrrrrrrrrrrrrBZbvji"
)

var acct *data.Account

func TestMain(m *testing.M) {
	var err error
	acct, err = data.NewAccountFromAddress(account)
	if err != nil {
		log.Fatalf("Invalid account %s: %v", account, err)
	}
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
		if err := got.New(tc.transaction, *acct); err != nil {
			t.Errorf("%v", err)
		}
		if got.String() != tc.exp {
			t.Errorf("expect csv %q, got %v", tc.exp, got)
		}
	}
}
