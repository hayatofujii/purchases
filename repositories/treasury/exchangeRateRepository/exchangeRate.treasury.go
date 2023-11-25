package exchangeRateTreasuryRepository

import (
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"time"

	exchangeRateData "haf.systems/purchases/models/exchange_rate"

	"haf.systems/purchases/utils"
)

const TREASURY_EXCHANGERATE_BASEURL = "https://api.fiscaldata.treasury.gov/services/api/fiscal_service/v1/accounting/od/rates_of_exchange"

type ExchangeRateFEDRepository struct {
	httpClient *http.Client
}

func NewExchangeRateTreasuryRepository(hc *http.Client) *ExchangeRateFEDRepository {
	return &ExchangeRateFEDRepository{
		httpClient: hc,
	}
}

type TreasuryExchangeRateDataEntry struct {
	ExchangeRate string `json:"exchange_rate"`
	RecordDate   string `json:"record_date"`
	Currency     string `json:"currency"`
}

type tresuryExchangeRateData struct {
	Data []TreasuryExchangeRateDataEntry `json:"data"`
}

func (r ExchangeRateFEDRepository) GetBestExchangeRate(currency string, date time.Time) (*exchangeRateData.ExchangeRateData, *utils.HTTPError) {

	field := "fields=exchange_rate,record_date,currency"
	sort := "sort=exchange_rate"

	filter := "filter="
	dateFilterGreater := "record_date:gte:" + date.AddDate(0, -6, 0).Format(time.DateOnly)
	dateFilterLesser := "record_date:lte:" + date.Format(time.DateOnly)
	countryCurrDesc := "country_currency_desc:in:" + currency

	url := TREASURY_EXCHANGERATE_BASEURL
	url += "?" + field
	url += "&" + sort
	url += "&" + filter + dateFilterGreater + "," + dateFilterLesser + "," + countryCurrDesc

	res, err := r.httpClient.Get(url)
	if err != nil {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("could not connect to upstream server: " + err.Error()),
		}
	}

	if res.StatusCode != http.StatusOK {
		return nil, &utils.HTTPError{
			StatusCode: res.StatusCode,
			Err:        fmt.Errorf("upstream response not OK: code %v", res.StatusCode),
		}
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("could not parse upstream data: " + err.Error()),
		}
	}

	var resParse tresuryExchangeRateData
	err = json.Unmarshal(resBody, &resParse)
	if err != nil {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("could not parse upstream data: " + err.Error()),
		}
	}

	if len(resParse.Data) == 0 {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusNotFound,
			Err:        fmt.Errorf("exchange rate for requested currency not found"),
		}
	}

	exRate := new(big.Rat)
	_, ratConvertOk := exRate.SetString(resParse.Data[0].ExchangeRate)

	if !ratConvertOk {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("could not parse upstream data: " + err.Error()),
		}
	}

	date, err = time.Parse(time.DateOnly, resParse.Data[0].RecordDate)
	if err != nil {
		return nil, &utils.HTTPError{
			StatusCode: http.StatusInternalServerError,
			Err:        fmt.Errorf("could not parse upstream data: " + err.Error()),
		}
	}

	return &exchangeRateData.ExchangeRateData{
		ExchangeRate: exRate,
		Date:         date,
		Currency:     resParse.Data[0].Currency,
	}, nil
}
