package purchaseFileRepository

// import (
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/suite"
// )

// is there a way to mock files?

// var allTestsFilter = func(_, _ string) (bool, error) { return true, nil }

// func TestPurchaseFileRepository(t *testing.T) {
// 	suite.Run(t, new(purchaseFileRepoTest))
// }

// type purchaseFileRepoTest struct {
// 	suite.Suite

// 	repo *PurchaseFileRepository
// }

// func (s *purchaseFileRepoTest) SetupSuite() {
// }

// func (s *purchaseFileRepoTest) SetupTest() {
// 	s.repo = NewPurchaseFileRepository("teste.db")
// }

// func (s *purchaseFileRepoTest) TestGetPurchase() {

// 	ok := testing.RunTests(
// 		allTestsFilter,
// 		[]testing.InternalTest{
// 			{
// 				Name: "Success",
// 				F: func(t *testing.T) {
// 				},
// 			},
// 			{
// 				Name: "ID not found",
// 				F: func(t *testing.T) {
// 				},
// 			},
// 		},
// 	)

// 	assert.Equal(s.T(), true, ok)
// }

// func (s *purchaseFileRepoTest) TestRecordPurchase() {

// 	ok := testing.RunTests(
// 		allTestsFilter,
// 		[]testing.InternalTest{
// 			{
// 				Name: "Success",
// 				F: func(t *testing.T) {
// 				},
// 			},
// 			{
// 				Name: "ID already exists",
// 				F: func(t *testing.T) {
// 				},
// 			},
// 			{
// 				Name: "fopen/append failure",
// 				F: func(t *testing.T) {
// 				},
// 			},
// 		},
// 	)

// 	assert.Equal(s.T(), true, ok)
// }
