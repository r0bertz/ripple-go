package csv

import (
	"fmt"

	"github.com/r0bertz/ripple/data"
)

// BitcoinTax represents bitcoin.tax csv format.
type BitcoinTax struct {
}

// Header prints bitcoin.tax csv header.
func (b BitcoinTax) Header() string {
	return "Date,Source,Action,Symbol,Volume,Currency,Price,Fee,FeeCurrency"
}

// Format prints Row in bitcoin.tax format.
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
func (b BitcoinTax) Format(r Row) (string, error) {
	t := r.TransactionWithMetaData
	var (
		Source      = "XRP Ledger"
		Action      Action
		Symbol      data.Currency
		Volume      data.Value
		Currency    data.Currency
		Price       data.Value
		Fee         data.Value
		FeeCurrency data.Currency
	)
	var (
		symbol data.Currency
		volume data.Value
	)
	if _, ok := r.m[xrp]; ok {
		symbol = xrp
		volume = r.m[symbol]
	} else if _, ok = r.m[usd]; ok {
		for k := range r.m {
			if !k.Equals(usd) {
				symbol = k
				break
			}
		}
		if symbol.Equals(cny) {
			return "", fmt.Errorf("not implemented: cny usd trade excluded, hash: %s", t.GetBase().Hash)
		}
		volume = r.m[symbol]
	} else {
		return "", fmt.Errorf("not implemented: no xrp or usd, hash: %s", t.GetBase().Hash)
	}
	Symbol = symbol
	if volume.IsNegative() {
		Action = SELL
		Volume = *volume.Negate()
	} else {
		Action = BUY
		Volume = volume
	}
	for k := range r.m {
		if !k.Equals(symbol) {
			Currency = k
			break
		}
	}
	ratio, err := r.m[Currency].Ratio(Volume)
	if err != nil {
		return "", fmt.Errorf("error calculating ratio: %v, hash: %s", err, t.GetBase().Hash)
	}
	Price = *ratio
	if Price.IsNegative() {
		Price = *Price.Negate()
	}
	Fee = t.GetBase().Fee
	FeeCurrency = xrp

	date := t.Date.Time()
	if Action == FEE {
		return fmt.Sprintf("%s,%s,%s,%s,,%s,,%.6f,", date.Format("2006-01-02 15:04:05 -0700"), Source, Action, Symbol, Currency, Fee.Float()), nil
	}
	return fmt.Sprintf("%s,%s,%s,%s,%.6f,%s,%.6f,%.6f,%s", date.Format("2006-01-02 15:04:05 -0700"), Source, Action, Symbol, Volume.Float(), Currency, Price.Float(), Fee.Float(), FeeCurrency), nil
}
