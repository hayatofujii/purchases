package purchaseService

import (
	"testing"

	exchangeRateRepoMock "haf.systems/purchases/mocks/repository/exchangeRate"
	purchaseRepoMock "haf.systems/purchases/mocks/repository/purchase"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// is there a way to mock files?

var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }

func TestPurchaseService(t *testing.T) {
	suite.Run(t, new(purchaseServiceTest))
}

type purchaseServiceTest struct {
	suite.Suite

	purchaseRepo     *purchaseRepoMock.PurchaseRepository
	exchangeRateRepo *exchangeRateRepoMock.ExchangeRateRepository

	s *PurchaseService
}

func (s *purchaseServiceTest) SetupSuite() {
}

func (s *purchaseServiceTest) SetupTest() {
	s.purchaseRepo = purchaseRepoMock.NewPurchaseRepository(s.T())
	s.exchangeRateRepo = exchangeRateRepoMock.NewExchangeRateRepository(s.T())
}

func (s *purchaseServiceTest) TestGetPurchase() {

	ok := testing.RunTests(
		allTestsFilter,
		[]testing.InternalTest{
			{
				Name: "Success",
				F: func(t *testing.T) {
				},
			},
			{
				Name: "ID not found",
				F: func(t *testing.T) {
				},
			},
		},
	)

	assert.Equal(s.T(), true, ok)
}

func (s *purchaseServiceTest) TestRecordPurchase() {
	ok := testing.RunTests(
		allTestsFilter,
		[]testing.InternalTest{
			{
				Name: "Success",
				F: func(t *testing.T) {
				},
			},
			{
				Name: "ID already exists",
				F: func(t *testing.T) {
				},
			},
			{
				Name: "fopen/append failure",
				F: func(t *testing.T) {
				},
			},
		},
	)

	assert.Equal(s.T(), true, ok)
}
