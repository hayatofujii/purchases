package purchaseService

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"

	exchangeRateModel "haf.systems/purchases/models/exchange_rate"
	purchaseModel "haf.systems/purchases/models/purchase"

	"haf.systems/purchases/utils"
)

type PurchaseRepository interface {
	GetPurchase(id string) (*purchaseModel.Purchase, *utils.HTTPError)
	RecordPurchase(id string, p purchaseModel.Purchase) error
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

func (p *PurchaseService) RegisterPurchase(id string, description string, date string, value float32) *utils.HTTPError {

	parsedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return &utils.HTTPError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("could not parse date"),
		}
	}

	purchase := purchaseModel.NewPurchase(description, parsedDate, value)
	err = p.purchaseRepository.RecordPurchase(id, purchase)
	if err != nil {
		return &utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	return nil
}

func (ps *PurchaseService) GetPurchase(id string) (*purchaseModel.Purchase, *utils.HTTPError) {

	p, e := ps.purchaseRepository.GetPurchase(id)
	if e != nil {
		return nil, e
	}

	return p, nil
}

func (ps *PurchaseService) GetConvertedPurchase(id string, currency string) (*purchaseModel.ConvertedPurchase, *utils.HTTPError) {

	p, e := ps.purchaseRepository.GetPurchase(id)
	if e != nil {
		return nil, e
	}

	rate, exchangeErr := ps.exchangeRateRepository.GetBestExchangeRate(currency, p.Date())
	if exchangeErr != nil {
		return nil, exchangeErr
	}

	return &purchaseModel.ConvertedPurchase{
		Purchase:       *p,
		ConvertedValue: p.Value() * rate.ExchangeRate,
		Currency:       rate.Currency,
		RateDate:       rate.Date,
	}, nil

}
