package controller

import (
	purchaseController "haf.systems/purchases/controllers/purchase"

	services "haf.systems/purchases/services"
)

type Controllers struct {
	PurchaseController *purchaseController.PurchaseController
}

func NewControllers(s *services.Services) *Controllers {
	return &Controllers{
		PurchaseController: purchaseController.NewPurchaseController(s.PurchaseService),
	}
}
