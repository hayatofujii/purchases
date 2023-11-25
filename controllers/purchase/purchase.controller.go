package purchaseController

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	purchaseModel "haf.systems/purchases/models/purchase"

	"haf.systems/purchases/utils"
)

type purchaseService interface {
	CreateID() string
	GetPurchase(string) (bool, *purchaseModel.Purchase)
	RegisterPurchase(string, string, string, string) (*bool, *utils.HTTPError)
	ConvertPurchase(*purchaseModel.Purchase, string) (*purchaseModel.ConvertedPurchase, *utils.HTTPError)
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

	alreadyRegistered, he := c.purchaseService.RegisterPurchase(id, req.Description, req.Date, fmt.Sprint(req.Value))
	if he != nil {
		ctx.Error(he)
		ctx.AbortWithStatusJSON(he.StatusCode, gin.H{
			"error": he.Error(),
		})
		return
	}

	if *alreadyRegistered {
		ctx.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PurchaseController) GetPurchase(ctx *gin.Context) {

	id := ctx.Params.ByName("id")

	exists, p := c.purchaseService.GetPurchase(id)

	if !exists {
		e := fmt.Errorf("not found")
		ctx.Error(e)
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	v, _ := strconv.ParseFloat(p.Value.FloatString(2), 64)
	response := gin.H{
		"id":          id,
		"value":       v,
		"description": p.Description,
		"date":        p.Date.Format(time.DateOnly),
	}

	currency := ctx.Query("currency")

	if currency != "" {
		convertions := map[string]gin.H{}

		currencies := strings.Split(currency, ",")

		for _, e := range currencies {
			if e != "" {
				conv, convErr := c.purchaseService.ConvertPurchase(p, e)

				if convErr != nil {
					convertions[e] = gin.H{
						"error": convErr.Error(),
					}
				} else {
					cv, _ := strconv.ParseFloat(conv.ConvertedValue.FloatString(2), 64)
					r, _ := strconv.ParseFloat(conv.Rate.FloatString(2), 64)

					convertions[e] = gin.H{
						"currency":        conv.Currency,
						"converted_value": cv,
						"rate":            r,
						"rate_date":       conv.RateDate.Format(time.DateOnly),
					}
				}
			}
		}

		response["convertions"] = convertions
	}

	ctx.IndentedJSON(http.StatusOK, response)

}
