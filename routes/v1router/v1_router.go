package v1router

import (
	controllers "haf.systems/purchases/controllers"

	v1PurchasesRoutes "haf.systems/purchases/routes/v1router/purchases"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(c *controllers.Controllers, router *gin.Engine) {
	v1 := router.Group("v1")

	v1PurchasesRoutes.SetupRoutes(c, v1)
}
