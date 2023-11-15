package purchaseModel

import (
	"time"
)

type Purchase struct {
	description string
	date        time.Time

	// big.Rat?
	value float32
}

func NewPurchase(_desc string, _date time.Time, _val float32) Purchase {
	return Purchase{
		description: _desc,
		date:        _date,
		value:       _val,
	}
}

func (a *Purchase) Description() string {
	return a.description
}

func (p *Purchase) Date() time.Time {
	return p.date
}

func (p *Purchase) Value() float32 {
	return p.value
}

type ConvertedPurchase struct {
	Purchase

	ConvertedValue float32
	Currency       string
	RateDate       time.Time
}
