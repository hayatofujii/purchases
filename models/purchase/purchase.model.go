package purchaseModel

import (
	"time"
)

type Purchase struct {
	Description string    `json:"description"`
	Date        time.Time `json:"date"`

	// big.Rat?
	Value float32 `json:"value"`
}

func NewPurchase(_desc string, _date time.Time, _val float32) Purchase {
	return Purchase{
		Description: _desc,
		Date:        _date,
		Value:       _val,
	}
}

type ConvertedPurchase struct {
	Purchase

	ConvertedValue float32   `json:"converted_value"`
	Currency       string    `json:"currency"`
	Rate           float32   `json:"rate"`
	RateDate       time.Time `json:"rate_date"`
}

type PurchaseSerial struct {
	Purchase
	ID string `json:"id"`
}

type ConvertedPurchaseSerial struct {
	ConvertedPurchase
	ID string `json:"id"`
}
