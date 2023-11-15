package services

import (
	// fedRepositories "haf.systems/purchases/repositories/fed"
	// fileRepositories "haf.systems/purchases/repositories/file"

	purchaseService "haf.systems/purchases/services/purchase"
)

type Services struct {
	PurchaseService *purchaseService.PurchaseService
}

func NewServices(pr purchaseService.PurchaseRepository, fxr purchaseService.ExchangeRateRepository) *Services {
	return &Services{
		PurchaseService: purchaseService.NewPurchaseService(pr, fxr),
	}
}
