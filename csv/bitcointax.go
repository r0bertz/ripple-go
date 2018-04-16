package csv

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/r0bertz/ripple/data"
)

// BitcoinTax represents bitcoin.tax csv format.
type BitcoinTax struct {
	Date        time.Time
	Source      string
	Action      Action
	Symbol      data.Currency
	Volume      data.Value
	Currency    data.Currency
	Price       data.Value
	Fee         data.Value
	FeeCurrency data.Currency
}

// New creates a Row from TransactionWithMetaData.
func (r *BitcoinTax) New(transaction, account string) error {
	var resp TxResponse
	dec := json.NewDecoder(strings.NewReader(transaction))
	if err := dec.Decode(&resp); err != nil {
		return fmt.Errorf("error decoding transaction: %v", err)
	}
	t := resp.Result
	switch t.GetTransactionType() {
	case data.ACCOUNT_SET, data.TRUST_SET, data.OFFER_CANCEL:
		// if err := accountRootBalanceChangeEqualsFee(t, account); err != nil {
		// 	return fmt.Errorf("not implemented: %v: hash %s", err, t.GetBase().Hash)
		// }
		// if t.GetBase().Account.String() != account {
		// 	return fmt.Errorf("not implemented: got account %s, expect %s: hash %s", t.GetBase().Account, account, t.GetBase().Hash)
		// }
		// r.Date = t.Date.Time()
		// r.Action = FEE
		// r.Symbol = xrp
		// r.Currency = xrp
		// r.Fee = t.GetBase().Fee
		// return nil
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
		if len(m) < 2 {
			// if err := accountRootBalanceChangeEqualsFee(t, account); err != nil {
			// 	return fmt.Errorf("not implemented: %v, hash: %s", err, t.GetBase().Hash)
			// }
			// if t.GetBase().Account.String() != account {
			// 	return fmt.Errorf("not implemented: got account %s, expect %s: %s", t.GetBase().Account, account, t.GetBase().Hash)
			// }
			// r.Date = t.Date.Time()
			// r.Action = FEE
			// r.Symbol = xrp
			// r.Currency = xrp
			// r.Fee = t.GetBase().Fee
			// return nil
			return fmt.Errorf("not implemented. fee. hash: %s", t.GetBase().Hash)
		}
		if len(m) != 2 {
			for k, v := range m {
				fmt.Printf("%s: %+v\n", k, v)
			}
			return fmt.Errorf("more than 2 balances, hash: %s", t.GetBase().Hash)
		}
		var (
			symbol data.Currency
			volume data.Value
			ok     bool
		)
		if _, ok = m[usd]; ok {
			for k := range m {
				if !k.Equals(usd) {
					symbol = k
					break
				}
			}
			if symbol.Equals(cny) {
				return fmt.Errorf("not implemented: cny usd trade excluded, hash: %s", t.GetBase().Hash)
			}
			volume = m[symbol]
		} else if _, ok = m[xrp]; ok {
			for k := range m {
				if !k.Equals(xrp) {
					symbol = k
					break
				}
			}
			volume = m[symbol]
		} else {
			return fmt.Errorf("not implemented: no xrp or usd, hash: %s", t.GetBase().Hash)
		}
		r.Date = t.Date.Time()
		r.Symbol = symbol
		if volume.IsNegative() {
			r.Action = SELL
			r.Volume = *volume.Negate()
		} else {
			r.Action = BUY
			r.Volume = volume
		}
		for k := range m {
			if !k.Equals(symbol) {
				r.Currency = k
				break
			}
		}
		ratio, err := m[r.Currency].Ratio(r.Volume)
		if err != nil {
			return fmt.Errorf("error calculating ratio: %v, hash: %s", err, t.GetBase().Hash)
		}
		r.Price = *ratio
		if r.Price.IsNegative() {
			r.Price = *r.Price.Negate()
		}
		r.Fee = t.GetBase().Fee
		r.FeeCurrency = xrp
		return nil
	}
	return fmt.Errorf("not implemented. hash: %s", t.GetBase().Hash)
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
func (r BitcoinTax) String() string {
	if r.Action == FEE {
		return fmt.Sprintf("%s,%s,%s,%s,,%s,,%.6f,", r.Date.Format("2006-01-02 15:04:05 -0700"), r.Source, r.Action, r.Symbol, r.Currency, r.Fee.Float())
	}
	return fmt.Sprintf("%s,%s,%s,%s,%.6f,%s,%.6f,%.6f,%s", r.Date.Format("2006-01-02 15:04:05 -0700"), r.Source, r.Action, r.Symbol, r.Volume.Float(), r.Currency, r.Price.Float(), r.Fee.Float(), r.FeeCurrency)
}

// DateTime returns Date.
func (r BitcoinTax) DateTime() time.Time {
	return r.Date
}
