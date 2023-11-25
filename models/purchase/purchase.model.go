package purchaseModel

import (
	"math/big"
	"time"
)

type Purchase struct {
	Description string    `json:"description"`
	Date        time.Time `json:"date"`

	Value *big.Rat
}

func NewPurchase(_desc string, _date time.Time, _val string) Purchase {
	r := new(big.Rat)
	r.SetString(_val)

	return Purchase{
		Description: _desc,
		Date:        _date,
		Value:       r,
	}
}

func (p Purchase) ValueFloat() float32 {
	f, _ := p.Value.Float32()
	return f
}

type PurchaseSerial struct {
	Purchase
	ID string `json:"id"`
}

type ConvertedPurchase struct {
	Purchase

	ConvertedValue *big.Rat
	Currency       string
	Rate           *big.Rat
	RateDate       time.Time
}

func (p ConvertedPurchase) ConvertedValueFloat() float32 {
	f, _ := p.ConvertedValue.Float32()
	return f
}

func (p ConvertedPurchase) RateFloat() float32 {
	f, _ := p.ConvertedValue.Float32()
	return f
}

type ConvertedPurchaseSerial struct {
	ConvertedPurchase
	ID string `json:"id"`
}
