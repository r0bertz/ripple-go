package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/r0bertz/ripple/data"
	"github.com/r0bertz/ripple/websockets"
)

var (
	addr    = flag.String("addr", "wss://s2.ripple.com:443", "wss service address")
	account = flag.String("account", "", "ripple account")
)

func main() {
	flag.Parse()

	if *account == "" {
		log.Fatal("--account not set")
	}
	acct, err := data.NewAccountFromAddress(*account)
	if err != nil {
		log.Fatal(err)
	}
	remote, err := websockets.NewRemote(*addr)
	if err != nil {
		log.Fatal(err)
	}

	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)

	for t := range remote.AccountTx(*acct, 10, -1, -1) {
		balances, err := t.Balances()
		if err != nil {
			log.Fatal(err)
		}
		if len(balances) == 0 {
			continue
		}
		m := map[data.Currency]data.Value{}
		b, ok := balances[*acct]
		if !ok {
			continue
		}
		for _, balance := range []data.Balance(*b) {
			if c, ok := m[balance.Currency]; ok {
				v, err := m[balance.Currency].Add(c)
				if err != nil {
					log.Fatal(err)
				}
				m[balance.Currency] = *v
				continue
			}
			m[balance.Currency] = balance.Change
		}
		var (
			dividend         data.Value
			dividendCurrency data.Currency
			divisor          data.Value
			divisorCurrency  data.Currency
		)
		for c, v := range m {
			if c.IsNative() {
				divisor = v
				divisorCurrency = c
				continue
			}
			dividend = v
			dividendCurrency = c
		}
		ratio, _ := data.NewValue("0", false)
		if len(m) == 2 {
			ratio, err = dividend.Ratio(divisor)
			if err != nil {
				log.Fatal(err)
			}
		}
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t%s\n", t.Date, dividendCurrency, dividend, divisorCurrency, divisor, ratio.Negate())
	}
	w.Flush()
}
