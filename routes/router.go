package routes

import (
	"github.com/gin-gonic/gin"

	"haf.systems/purchases/routes/v1router"

	controllers "haf.systems/purchases/controllers"
)

func SetRoutes(c *controllers.Controllers, router *gin.Engine) {
	v1router.SetupRoutes(c, router)
}
