// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package exchange_repository_mock

import (
	mock "github.com/stretchr/testify/mock"
	exchangeRateDataModel "haf.systems/purchases/models/exchange_rate"

	time "time"

	utils "haf.systems/purchases/utils"
)

// ExchangeRateRepository is an autogenerated mock type for the ExchangeRateRepository type
type ExchangeRateRepository struct {
	mock.Mock
}

// GetBestExchangeRate provides a mock function with given fields: currency, date
func (_m *ExchangeRateRepository) GetBestExchangeRate(currency string, date time.Time) (*exchangeRateDataModel.ExchangeRateData, *utils.HTTPError) {
	ret := _m.Called(currency, date)

	var r0 *exchangeRateDataModel.ExchangeRateData
	var r1 *utils.HTTPError
	if rf, ok := ret.Get(0).(func(string, time.Time) (*exchangeRateDataModel.ExchangeRateData, *utils.HTTPError)); ok {
		return rf(currency, date)
	}
	if rf, ok := ret.Get(0).(func(string, time.Time) *exchangeRateDataModel.ExchangeRateData); ok {
		r0 = rf(currency, date)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*exchangeRateDataModel.ExchangeRateData)
		}
	}

	if rf, ok := ret.Get(1).(func(string, time.Time) *utils.HTTPError); ok {
		r1 = rf(currency, date)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*utils.HTTPError)
		}
	}

	return r0, r1
}

// NewExchangeRateRepository creates a new instance of ExchangeRateRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewExchangeRateRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *ExchangeRateRepository {
	mock := &ExchangeRateRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}