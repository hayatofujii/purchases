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

func (p *PurchaseService) RegisterPurchase(id string, description string, date string, value float32) (*bool, *utils.HTTPError) {

	parsedDate, err := time.Parse(time.DateOnly, date)
	if err != nil {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusBadRequest,
			Err:        fmt.Errorf("could not parse date"),
		}
	}

	purchase := purchaseModel.NewPurchase(description, parsedDate, value)
	ok, err := p.purchaseRepository.RecordPurchase(id, purchase)

	if err != nil {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        err,
		}
	}

	retExist := true

	if !ok {
		return &retExist, nil
	}

	return &retExist, nil
}

func (ps *PurchaseService) GetPurchase(id string) (bool, *purchaseModel.PurchaseSerial) {
	e, p := ps.purchaseRepository.GetPurchase(id)

	if !e {
		return e, nil
	}

	return e, &purchaseModel.PurchaseSerial{
		ID:       id,
		Purchase: *p,
	}
}

func (ps *PurchaseService) GetConvertedPurchase(id string, currency string) (bool, *purchaseModel.ConvertedPurchaseSerial, *utils.HTTPError) {

	ok, p := ps.purchaseRepository.GetPurchase(id)

	if !ok {
		return false, nil, nil
	}

	rate, exchangeErr := ps.exchangeRateRepository.GetBestExchangeRate(currency, p.Date)
	if exchangeErr != nil {
		return true, nil, exchangeErr
	}

	return ok, &purchaseModel.ConvertedPurchaseSerial{
		ConvertedPurchase: purchaseModel.ConvertedPurchase{
			Purchase: purchaseModel.Purchase{
				Description: p.Description,
				Date:        p.Date,
				Value:       p.Value,
			},
			ConvertedValue: p.Value * rate.ExchangeRate,
			Currency:       rate.Currency,
			Rate:           rate.ExchangeRate,
			RateDate:       rate.Date,
		},
		ID: id,
	}, nil
}
