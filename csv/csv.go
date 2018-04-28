package csv

import (
	"container/heap"
	"fmt"
	"reflect"
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
	FormatterFactory = map[string]func() Formatter{
		"bitcoin.tax":    func() Formatter { return &BitcoinTax{} },
		"cointracker.io": func() Formatter { return &CoinTrackerIO{} },
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

// Row contains fields common to all CSV formats.
type Row struct {
	data.TransactionWithMetaData
	// m maps Currency to changed value.
	m map[data.Currency]data.Value
}

// URL returns the URL of the transaction that's associated with this Row.
func (r Row) URL() string {
	return fmt.Sprintf("http://ripplescan.com/transactions/%s", r.TransactionWithMetaData.GetBase().Hash)
}

// DateTime returns Date
func (r Row) DateTime() time.Time {
	return r.TransactionWithMetaData.Date.Time()
}

// Heap is a max heap of Rows
type Heap []Row

// Len returns length of Heap.
func (h Heap) Len() int { return len(h) }

// Swap swaps elements at i and j.
func (h Heap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

// Less returns true if element j have a smaller timestamp than element i.
func (h Heap) Less(i, j int) bool { return h[j].DateTime().Before(h[i].DateTime()) }

// Push pushes x on to heap.
func (h *Heap) Push(x interface{}) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Row))
}

// Pop pops the last row off heap.
func (h *Heap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

// CSV represents all rows in a csv file.
type CSV struct {
	Rows    Heap
	Account data.Account
	Related []data.Account
}

// New returns a new CSV.
func New(account data.Account, related []data.Account) *CSV {
	rv := &CSV{
		Rows:    Heap{},
		Account: account,
		Related: related,
	}
	heap.Init(&rv.Rows)
	return rv
}

// Add adds a new transaction.
func (c *CSV) Add(t data.TransactionWithMetaData) error {
	switch t.GetTransactionType() {
	case data.ACCOUNT_SET, data.TRUST_SET, data.OFFER_CANCEL:
		return fmt.Errorf("not implemented. fee. hash: %s", t.GetBase().Hash)
	case data.PAYMENT, data.OFFER_CREATE:
		bm, err := t.Balances()
		if err != nil {
			return fmt.Errorf("error getting balances: %v, hash: %s", err, t.GetBase().Hash)
		}
		b, ok := bm[c.Account]
		if !ok {
			return fmt.Errorf("not implemented. fee. hash: %s", t.GetBase().Hash)
		}
		m := map[data.Currency]data.Value{}
		for _, b := range []data.Balance(*b) {
			c, ok := m[b.Currency]
			if !ok {
				m[b.Currency] = b.Change
				continue

			}
			v, err := c.Add(b.Change)
			if err != nil {
				return fmt.Errorf("error adding balance changes %s and %s: %v, hash: %s", c, b.Change, err, t.GetBase().Hash)
			}
			m[b.Currency] = *v
		}
		if len(m) == 1 {
			v, ok := m[xrp]
			if !ok {
				keys := reflect.ValueOf(m).MapKeys()
				currency := keys[0].Interface().(data.Currency)
				return fmt.Errorf("not implemented. depositing IOU %s %s. hash: %s", currency, m[currency], t.GetBase().Hash)
			}
			p := t.Transaction.(*data.Payment)
			for _, a := range c.Related {
				if p.Account.Equals(a) {
					return fmt.Errorf("not implemented. payment from related account %s: %s. hash: %s", a, v, t.GetBase().Hash)
				}
				if p.Destination.Equals(a) {
					return fmt.Errorf("not implemented. payment to related account %s: %s. hash: %s", a, v, t.GetBase().Hash)
				}
			}
			if v.IsNegative() {
				return fmt.Errorf("not implemented. sent out xrp %s, hash: %s", v, t.GetBase().Hash)
			}
		}
		if len(m) > 2 {
			return fmt.Errorf("more than 2 currencies: %+v, hash: %s", m, t.GetBase().Hash)
		}
		r := Row{TransactionWithMetaData: t, m: m}
		heap.Push(&c.Rows, r)
		return nil
	}
	return fmt.Errorf("not implemented. hash: %s", t.GetBase().Hash)
}

// Formatter prints csv file header and rows.
type Formatter interface {
	Header() string
	Format(r Row) (string, error)
}
