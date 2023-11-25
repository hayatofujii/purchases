package exchangeRateDataModel

import (
	"math/big"
	"time"
)

type ExchangeRateData struct {
	ExchangeRate *big.Rat
	Date         time.Time
	Currency     string
}
