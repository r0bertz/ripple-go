package csv

import (
	"fmt"

	"github.com/r0bertz/ripple/data"
)

// CoinTrackerIO represents cointracker.io csv format.
type CoinTrackerIO struct {
}

// Header prints cointracker.io csv header.
func (c CoinTrackerIO) Header() string {
	return "Date,Received Quantity,Currency,Sent Quantity,Currency"
}

// Format prints Row in cointracker.io csv format.
// The return value contains the following columns in this order:
//   * Date (date and time as MM/DD/YYYY HH:mm:ss)
//   * Received Quantity
//   * Currency (specify currency such as USD, GBP, EUR or coins, BTC or LTC)
//   * Sent Quantity
//   * Currency (specify currency such as USD, GBP, EUR or coins, BTC or LTC)
func (c CoinTrackerIO) Format(r Row) (string, error) {
	var (
		ReceivedCurrency data.Currency
		Received         data.Value
		SentCurrency     data.Currency
		Sent             data.Value
	)
	for c, q := range r.m {
		if q.IsNegative() {
			SentCurrency = c
			Sent = *q.Negate()
		} else {
			ReceivedCurrency = c
			Received = q
		}
	}

	date := r.TransactionWithMetaData.Date.Time()
	if Received.IsZero() {
		return fmt.Sprintf("%s,,,%.6f,%s", date.Format("01/02/2006 15:04:05"), Sent.Float(), SentCurrency), nil
	}
	if Sent.IsZero() {
		return fmt.Sprintf("%s,%.6f,%s,,", date.Format("01/02/2006 15:04:05"), Received.Float(), ReceivedCurrency), nil
	}
	return fmt.Sprintf("%s,%.6f,%s,%.6f,%s", date.Format("01/02/2006 15:04:05"), Received.Float(), ReceivedCurrency, Sent.Float(), SentCurrency), nil
}
