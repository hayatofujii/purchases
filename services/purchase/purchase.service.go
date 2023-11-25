package purchaseService

import (
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/google/uuid"

	exchangeRateModel "haf.systems/purchases/models/exchange_rate"
	purchaseModel "haf.systems/purchases/models/purchase"

	"haf.systems/purchases/utils"
)

type PurchaseRepository interface {
	GetPurchase(id string) (bool, *purchaseModel.Purchase)
	RecordPurchase(id string, p purchaseModel.Purchase) (bool, error)
}

type ExchangeRateRepository interface {
	GetBestExchangeRate(currency string, date time.Time) (*exchangeRateModel.ExchangeRateData, *utils.HTTPError)
}

type PurchaseService struct {
	purchaseRepository     PurchaseRepository
	exchangeRateRepository ExchangeRateRepository
}

func NewPurchaseService(pr PurchaseRepository, er ExchangeRateRepository) *PurchaseService {
	return &PurchaseService{
		purchaseRepository:     pr,
		exchangeRateRepository: er,
	}
}

func (p *PurchaseService) CreateID() string {
	return uuid.New().String()
}

func (p *PurchaseService) RegisterPurchase(id string, description string, date string, value string) (*bool, *utils.HTTPError) {

	parsedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("could not parse date"),
		}
	}

	purchase := purchaseModel.NewPurchase(description, parsedDate, value)
	registered, err := p.purchaseRepository.RecordPurchase(id, purchase)

	if err != nil {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	retExist := !registered

	return &retExist, nil
}

func (ps *PurchaseService) GetPurchase(id string) (bool, *purchaseModel.Purchase) {
	e, p := ps.purchaseRepository.GetPurchase(id)

	if !e {
		return e, nil
	}

	return e, p
}

func (ps *PurchaseService) ConvertPurchase(p *purchaseModel.Purchase, currency string) (*purchaseModel.ConvertedPurchase, *utils.HTTPError) {
	rate, exchangeErr := ps.exchangeRateRepository.GetBestExchangeRate(currency, p.Date)
	if exchangeErr != nil {
		return nil, exchangeErr
	}

	mul := new(big.Rat)
	mul = mul.Mul(p.Value, rate.ExchangeRate)

	return &purchaseModel.ConvertedPurchase{
		ConvertedValue: mul,
		Currency:       rate.Currency,
		Rate:           rate.ExchangeRate,
		RateDate:       rate.Date,
	}, nil
}
