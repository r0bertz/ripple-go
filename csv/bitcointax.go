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
		FeeCurrency = xrp
		Fee         = t.GetBase().Fee
		Action      Action
		Symbol      *data.Currency
		Volume      *data.Value
		Currency    *data.Currency
		Price       *data.Value
		err         error
	)
	if v, ok := r.m[xrp]; ok {
		Symbol = &xrp
		Volume = &v
	} else if _, ok = r.m[usd]; ok {
		for k, v := range r.m {
			if !k.Equals(usd) {
				Symbol = &k
				Volume = &v
				break
			}
		}
		if Symbol == nil {
			return "", fmt.Errorf("usd deposit or withdrawal, hash: %s", t.GetBase().Hash)
		}
		if Symbol.Equals(cny) {
			return "", fmt.Errorf("not implemented: cny usd trade excluded, hash: %s", t.GetBase().Hash)
		}
	} else {
		return "", fmt.Errorf("not implemented: no xrp or usd, hash: %s", t.GetBase().Hash)
	}
	if len(r.m) == 0 {
		Action = FEE
	} else if Volume.IsNegative() {
		Action = SELL
		Volume = Volume.Negate()
	} else {
		Action = BUY
	}
	for k := range r.m {
		if !k.Equals(*Symbol) {
			Currency = &k
			break
		}
	}
	if Currency != nil {
		Price, err = r.m[*Currency].Ratio(*Volume)
		if err != nil {
			return "", fmt.Errorf("error calculating ratio: %v, hash: %s", err, t.GetBase().Hash)
		}
		if Price.IsNegative() {
			Price = Price.Negate()
		}
	}

	date := t.Date.Time()
	if Action == FEE {
		return fmt.Sprintf("%s,%s,%s,%s,,,,%.6f,%s", date.Format("2006-01-02 15:04:05 -0700"), Source, Action, Symbol, Fee.Float(), FeeCurrency), nil
	}
	if Currency == nil {
		return fmt.Sprintf("%s,%s,%s,%s,%.6f,XRP,0.000000,%.6f,%s", date.Format("2006-01-02 15:04:05 -0700"), Source, Action, Symbol, Volume.Float(), Fee.Float(), FeeCurrency), nil
	}
	return fmt.Sprintf("%s,%s,%s,%s,%.6f,%s,%.6f,%.6f,%s", date.Format("2006-01-02 15:04:05 -0700"), Source, Action, Symbol, Volume.Float(), Currency, Price.Float(), Fee.Float(), FeeCurrency), nil
}
