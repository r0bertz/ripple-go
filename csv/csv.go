package csv

import (
	"errors"
	"fmt"
	"time"

	"github.com/r0bertz/ripple/data"
	"github.com/r0bertz/ripple/websockets"
)

const (
	// UNKNOWN is the default value for action.
	UNKNOWN Action = iota
	// BUY means coins are bought in the transaction.
	BUY
	// SELL means coins are sold in the transaction.
	SELL
	// FEE means fees are charged in the transaction.
	FEE
)

// Factory returns a function that return a Row in given Row type.
var (
	Factory = map[string]func() Row{
		"bitcointax":     func() Row { return &BitcoinTax{} },
		"cointracker.io": func() Row { return &CoinTrackerIO{} },
	}
	xrp, _ = data.NewCurrency("XRP")
	usd, _ = data.NewCurrency("USD")
	cny, _ = data.NewCurrency("CNY")
)

// Action is one of BUY, SELL or FEE.
type Action int

func (a Action) String() string {
	names := []string{
		"UNKNOWN",
		"BUY",
		"SELL",
		"FEE",
	}
	if a > FEE || a < UNKNOWN {
		return "invalid action"
	}
	return names[a]
}

// TxResponse is the response of rippled tx method. https://ripple.com/build/rippled-apis/#tx
type TxResponse struct {
	Result websockets.TxResult
	Status string
	Type   string
}

// Base contains fields common to all CSV formats.
type Base struct {
	TxResult websockets.TxResult
}

// TxURL returns the URL of the transaction that's associated with this Row.
func (b Base) TxURL() string {
	return fmt.Sprintf("http://ripplescan.com/transactions/%s", b.TxResult.GetBase().Hash)
}

// DateTime returns Date
func (b Base) DateTime() time.Time {
	return b.TxResult.Date.Time()
}

// Row represents one row in csv.
type Row interface {
	New(transaction string, account data.Account, related []data.Account) error
	String() string
	TxURL() string
	DateTime() time.Time
}

func accountRootBalanceChangeEqualsFee(t websockets.TxResult, account string) error {
	for _, n := range t.MetaData.AffectedNodes {
		node, final, previous, state := n.AffectedNode()
		if state == data.Modified && node.LedgerEntryType == data.ACCOUNT_ROOT {
			f := final.(*data.AccountRoot)
			p := previous.(*data.AccountRoot)
			if f != nil && p != nil && f.Balance != nil && p.Balance != nil && f.Account != nil && f.Account.String() == account {
				diff, err := p.Balance.Subtract(*f.Balance)
				if err != nil {
					return err
				}
				if diff.Equals(t.GetBase().Fee) {
					return nil
				}
				// Payment to account, etc.
				return errors.New("account root balance change not equals to fee")
			}
		}
	}
	// Owner count change, etc.
	return errors.New("no account root blance change")
}

// Slice is a slice of Rows.
type Slice []Row

// Len returns length of Slice.
func (s Slice) Len() int { return len(s) }

// Swap swaps elements at i and j.
func (s Slice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less returns true if element i have a smaller timestamp than element j.
func (s Slice) Less(i, j int) bool { return s[i].DateTime().Before(s[j].DateTime()) }
