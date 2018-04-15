package csv

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
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
	timeFormat = "2006-01-02 15:04:05 -0700"
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

// Response is the response of rippled tx method. https://ripple.com/build/rippled-apis/#tx
type Response struct {
	Result websockets.TxResult
	Status string
	Type   string
}

// Row represents one row in csv.
type Row struct {
	Date        time.Time
	Source      string
	Action      Action
	Symbol      string
	Volume      float64
	Currency    string
	Price       float64
	Fee         data.Value
	FeeCurrency string
}

// NewRow creates a Row from TransactionWithMetaData.
func NewRow(transaction, account string) (Row, error) {
	var (
		rv Row
		r  Response
	)
	dec := json.NewDecoder(strings.NewReader(transaction))
	if err := dec.Decode(&r); err != nil {
		log.Printf("error decoding response: %v", err)
		if e, ok := err.(*json.SyntaxError); ok {
			log.Printf("syntax error at byte offset %d", e.Offset)
		}
		log.Printf("transaction: %q", transaction)
		return rv, err
	}
	t := r.Result
	b, _ := json.MarshalIndent(t, "", "  ")
	switch t.GetTransactionType() {
	case data.ACCOUNT_SET, data.TRUST_SET, data.OFFER_CANCEL:
		if t.GetBase().Account.String() != account {
			return rv, fmt.Errorf("got account %s, expect %s", t.GetBase().Account, account)
		}
		if err := accountRootBalanceChangeEqualsFee(t, account); err != nil {
			return rv, fmt.Errorf("%v: %s", err, string(b))
		}
		rv.Date = t.Date.Time()
		rv.Action = FEE
		rv.Symbol = "XRP"
		rv.Currency = "XRP"
		rv.Fee = t.GetBase().Fee
		return rv, nil
	default:
		balances, err := t.Balances()
		if err != nil {
			return rv, err
		}
		m := map[data.Currency]struct {
			Balance data.Value
			Change  data.Value
		}{}
		for _, b := range balances {
			if b.Account.String() == account {
				m[b.Currency] = struct {
					Balance data.Value
					Change  data.Value
				}{
					b.Balance,
					b.Change,
				}
			}
		}
		if len(m) < 2 {
			if err := accountRootBalanceChangeEqualsFee(t, account); err == nil {
				rv.Date = t.Date.Time()
				rv.Action = FEE
				rv.Symbol = "XRP"
				rv.Currency = "XRP"
				rv.Fee = t.GetBase().Fee
				return rv, nil
			}
			// account receives payment, etc.
		} else {
			if len(m) != 2 {
				for k, v := range m {
					fmt.Printf("%s: %+v\n", k, v)
				}
				fmt.Printf("accountBalances: %v, hash: %s\n", len(m), t.GetBase().Hash)
			} else {
				// TODO
			}
		}
	}
	return rv, errors.New("not implemented")
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
				return errors.New("account root balance change not equals to fee")
			}
		}
	}
	return errors.New("no account root blance change")
}

// The return value contains the following columns in this order:
//   * Date (date and time as YYYY-MM-DD HH:mm:ss Z)
//   * Source (optional, such as an exchange name like MtGox or gift, donation, etc)
//   * Action (BUY, SELL or FEE)
//   * Symbol (XRP)
//   * Volume (number of coins traded - ignore if FEE)
//   * Currency (specify currency such as USD, GBP, EUR or coins, BTC or LTC)
//   * Price (price per coin in Currency or blank for lookup - ignore if FEE)
//   * Fee (any additional costs of the trade)
//   * FeeCurrency (currency of fee if different than Currency)
func (r Row) String() string {
	if r.Action == FEE {
		return fmt.Sprintf("%s,%s,%s,%s,,%s,,%.6f,", r.Date.Format(timeFormat), r.Source, r.Action, r.Symbol, r.Currency, r.Fee.Float())
	}
	return fmt.Sprintf("%s,%s,%s,%s,%.6f,%s,%.6f,%.6f,%s", r.Date.Format(timeFormat), r.Source, r.Action, r.Symbol, r.Volume, r.Currency, r.Price, r.Fee, r.FeeCurrency)
}

// Slice is a slice of Rows.
type Slice []Row

// Len returns length of Slice.
func (s Slice) Len() int { return len(s) }

// Swap swaps elements at i and j.
func (s Slice) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

// Less returns true if element i have a smaller timestamp than element j.
func (s Slice) Less(i, j int) bool { return s[i].Date.Before(s[j].Date) }
