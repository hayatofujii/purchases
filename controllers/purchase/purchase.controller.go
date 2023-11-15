package purchaseController

import (
	"net/http"

	"github.com/gin-gonic/gin"

	purchaseModel "haf.systems/purchases/models/purchase"

	"haf.systems/purchases/utils"
)

type purchaseService interface {
	CreateID() string
	RegisterPurchase(string, string, string, float32) *utils.HTTPError
	GetPurchase(string) (*purchaseModel.Purchase, *utils.HTTPError)
	GetConvertedPurchase(string, string) (*purchaseModel.ConvertedPurchase, *utils.HTTPError)
}

type PurchaseController struct {
	purchaseService purchaseService
}

func NewPurchaseController(s purchaseService) *PurchaseController {
	return &PurchaseController{
		purchaseService: s,
	}
}

func (c *PurchaseController) RequestUUID(ctx *gin.Context) {
	ctx.IndentedJSON(http.StatusOK, gin.H{
		"id": c.purchaseService.CreateID(),
	})
}

const (
	PURCHASE_DESCRIPTION_MAX_LEN = 50
)

func (c *PurchaseController) RegisterPurchase(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	var req struct {
		Description string  `json:"description"`
		Amount      float32 `json:"amount"`
		Date        string  `json:"date"`
	}

	e := ctx.BindJSON(&req)
	if e != nil {
		ctx.Error(e)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": e.Error(),
		})
		return
	}

	if len(req.Description) > PURCHASE_DESCRIPTION_MAX_LEN {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "purchase description cannot be longer than 50 characters",
		})
		return
	}

	he := c.purchaseService.RegisterPurchase(id, req.Description, req.Date, req.Amount)
	if he != nil {
		ctx.Error(he)
		ctx.AbortWithStatusJSON(he.StatusCode, gin.H{
			"error": he.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, gin.H{
		"status": "ok",
	})
}

func (c *PurchaseController) GetPurchase(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	queryParams := ctx.Request.URL.Query()

	// let's get the first currency for the time being
	currency := queryParams["currency"][0]

	var p any
	var e *utils.HTTPError

	if currency != "" {
		p, e = c.purchaseService.GetConvertedPurchase(id, currency)
	} else {
		p, e = c.purchaseService.GetPurchase(id)
	}

	if e != nil {
		ctx.Error(e)
		ctx.AbortWithStatusJSON(e.StatusCode, gin.H{
			"error": e.Error(),
		})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, p)

}
