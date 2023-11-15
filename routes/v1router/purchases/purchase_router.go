package purchaseRoutes

import (
	"github.com/gin-gonic/gin"
	controllers "haf.systems/purchases/controllers"
)

func SetupRoutes(c *controllers.Controllers, parentGroup *gin.RouterGroup) {
	purchaseGroup := parentGroup.Group("purchase")

	purchaseGroup.POST("/", c.PurchaseController.RequestUUID)
	purchaseGroup.PUT("/:id", c.PurchaseController.RegisterPurchase)
	purchaseGroup.GET("/:id", c.PurchaseController.GetPurchase)
}
