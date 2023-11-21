package exchangeRateDataModel

import "time"

type ExchangeRateData struct {
	ExchangeRate float32
	Date         time.Time
	Currency     string
}
