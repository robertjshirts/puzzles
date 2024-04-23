package handler

import (
	"github.com/puzzles/services/orders/dal"
	"github.com/puzzles/services/orders/gen"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	db *dal.SQLDal
}

func NewOrderHandler(db *dal.SQLDal) *OrderHandler {
	return &OrderHandler{
		db: db,
	}
}

func (h *OrderHandler) CreateNewOrder(c *gin.Context) {
	var orderJSON gen.NewOrderInfo
	if err := c.BindJSON(&orderJSON); err != nil {
		c.JSON(400, gen.Error{Code: 400, Message: "invalid request"})
		return
	}

	orderInfo, shippingInfo, paymentInfo, puzzles, genErr := dal.ToDALModels(&orderJSON)
	if genErr != nil {
		c.JSON(genErr.Code, genErr)
		return
	}

	genErr = h.db.CreateOrder(*orderInfo, *paymentInfo, *shippingInfo, *puzzles)
	if genErr != nil {
		c.JSON(genErr.Code, genErr)
		return
	}

	apiOrderModel := dal.ToApiModel(*orderInfo, *paymentInfo, *shippingInfo, *puzzles)
	if apiOrderModel == nil {
		c.JSON(500, gen.Error{Code: 500, Message: "There was an issue converting from DAL models to API Models"})
		return
	}

	c.JSON(201, apiOrderModel)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context, id string) {
	genErr := h.db.DeleteOrder(id)
	if genErr != nil {
		c.JSON(genErr.Code, genErr)
		return
	}
	c.Status(204)
}

func (h *OrderHandler) GetOrder(c *gin.Context, id string) {
	orderInfo, payment, shipping, puzzles, genErr := h.db.GetOrder(id)
	if genErr != nil {
		c.JSON(genErr.Code, genErr)
		return
	}

	apiOrderModel := dal.ToApiModel(*orderInfo, *payment, *shipping, *puzzles)
	c.JSON(200, *apiOrderModel)
}
