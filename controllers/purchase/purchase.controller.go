package purchaseController

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	purchaseModel "haf.systems/purchases/models/purchase"

	"haf.systems/purchases/utils"
)

type purchaseService interface {
	CreateID() string
	GetPurchase(string) (bool, *purchaseModel.PurchaseSerial)
	RegisterPurchase(string, string, string, string) (*bool, *utils.HTTPError)
	GetConvertedPurchase(string, string) (bool, *purchaseModel.ConvertedPurchaseSerial, *utils.HTTPError)
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
		Value       float32 `json:"value"`
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

	registered, he := c.purchaseService.RegisterPurchase(id, req.Description, req.Date, fmt.Sprint(req.Value))
	if he != nil {
		ctx.Error(he)
		ctx.AbortWithStatusJSON(he.StatusCode, gin.H{
			"error": he.Error(),
		})
		return
	}

	if !*registered {
		ctx.IndentedJSON(http.StatusNoContent, gin.H{})
		return
	}

	ctx.IndentedJSON(http.StatusNoContent, gin.H{})
}

func (c *PurchaseController) GetPurchase(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	// let's get the first currency for the time being
	currency := ctx.Query("currency")

	if currency != "" {
		exists, p, e := c.purchaseService.GetConvertedPurchase(id, currency)

		if e != nil {
			ctx.Error(e)
			ctx.AbortWithStatusJSON(e.StatusCode, gin.H{
				"error": e.Error(),
			})
			return
		}

		if !exists {
			e := fmt.Errorf("not found")
			ctx.Error(e)
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
			return
		}

		v, _ := strconv.ParseFloat(p.Value.FloatString(2), 64)
		cv, _ := strconv.ParseFloat(p.ConvertedValue.FloatString(2), 64)
		r, _ := strconv.ParseFloat(p.Rate.FloatString(2), 64)

		ctx.IndentedJSON(http.StatusOK, gin.H{
			"id":          p.ID,
			"value":       v,
			"description": p.Description,
			"date":        p.Date.Format(time.DateOnly),

			"converted_value:": cv,
			"currency":         p.Currency,
			"rate":             r,
			"rate_date":        p.RateDate.Format(time.DateOnly),
		})

	} else {
		exists, p := c.purchaseService.GetPurchase(id)

		if !exists {
			e := fmt.Errorf("not found")
			ctx.Error(e)
			ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{})
			return
		}

		v, _ := strconv.ParseFloat(p.Value.FloatString(2), 64)

		ctx.IndentedJSON(http.StatusOK, gin.H{
			"id":          p.ID,
			"value":       v,
			"description": p.Description,
			"date":        p.Date.Format(time.DateOnly),
		})
	}

}
