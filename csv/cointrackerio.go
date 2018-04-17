package csv

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/r0bertz/ripple/data"
)

// CoinTrackerIO represents cointracker.io csv format.
type CoinTrackerIO struct {
	Base
	Received         data.Value
	ReceivedCurrency data.Currency
	Sent             data.Value
	SentCurrency     data.Currency
}

// New creates a Row from TransactionWithMetaData.
func (r *CoinTrackerIO) New(transaction, account string) error {
	var resp TxResponse
	dec := json.NewDecoder(strings.NewReader(transaction))
	if err := dec.Decode(&resp); err != nil {
		return fmt.Errorf("error decoding transaction: %v", err)
	}
	t := resp.Result
	switch t.GetTransactionType() {
	case data.ACCOUNT_SET, data.TRUST_SET, data.OFFER_CANCEL:
		return fmt.Errorf("not implemented. fee. hash: %s", t.GetBase().Hash)
	case data.PAYMENT, data.OFFER_CREATE:
		balances, err := t.Balances()
		if err != nil {
			return fmt.Errorf("error getting balances: %v, hash: %s", err, t.GetBase().Hash)
		}
		m := map[data.Currency]data.Value{}
		for _, b := range balances {
			if b.Account.String() == account {
				m[b.Currency] = b.Change
			}
		}
		if len(m) == 0 {
			return fmt.Errorf("not implemented: no balance, %s, hash: %s", t.Date, t.GetBase().Hash)
		}
		if len(m) > 2 {
			for k, v := range m {
				fmt.Printf("%s: %+v\n", k, v)
			}
			return fmt.Errorf("more than 2 balances, hash: %s", t.GetBase().Hash)
		}
		r.Date = t.Date.Time()
		r.Hash = t.GetBase().Hash
		for c, q := range m {
			if q.IsNegative() {
				r.SentCurrency = c
				r.Sent = *q.Negate()
			} else {
				r.ReceivedCurrency = c
				r.Received = q
			}
		}
		return nil
	}
	return fmt.Errorf("not implemented. hash: %s", t.GetBase().Hash)
}

// The return value contains the following columns in this order:
//   * Date (date and time as MM/DD/YYYY HH:mm:ss)
//   * Received Quantity
//   * Currency (specify currency such as USD, GBP, EUR or coins, BTC or LTC)
//   * Sent Quantity
//   * Currency (specify currency such as USD, GBP, EUR or coins, BTC or LTC)
func (r CoinTrackerIO) String() string {
	if r.Received.IsZero() {
		return fmt.Sprintf("%s,,,%.6f,%s", r.Date.Format("01/02/2006 15:04:05"), r.Sent.Float(), r.SentCurrency)
	}
	if r.Sent.IsZero() {
		return fmt.Sprintf("%s,%.6f,%s,,", r.Date.Format("01/02/2006 15:04:05"), r.Received.Float(), r.ReceivedCurrency)
	}
	return fmt.Sprintf("%s,%.6f,%s,%.6f,%s", r.Date.Format("01/02/2006 15:04:05"), r.Received.Float(), r.ReceivedCurrency, r.Sent.Float(), r.SentCurrency)
}
