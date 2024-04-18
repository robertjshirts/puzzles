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

	orderInfo, shippingInfo, paymentInfo, puzzles := dal.ToDALModels(&orderJSON)
	if orderInfo == nil || shippingInfo == nil || paymentInfo == nil || puzzles == nil {
		c.JSON(400, gen.Error{Code: 400, Message: "invalid request"})
		return
	}

	err := h.db.CreateOrderInfo(*orderInfo, *paymentInfo, *shippingInfo, *puzzles)
	if err != nil {
		c.JSON(500, gen.Error{Code: 500, Message: err.Error()})
		return
	}

	apiOrderModel := dal.ToApiModel(*orderInfo, *shippingInfo, *paymentInfo, *puzzles)
	if apiOrderModel == nil {
		c.JSON(500, gen.Error{Code: 500, Message: "There was an issue converting from DAL models to API Models"})
		return
	}

	c.JSON(201, apiOrderModel)
}

func (h *OrderHandler) DeleteOrder(c *gin.Context, id string) {
	c.Status(204)
}

func (h *OrderHandler) GetOrder(c *gin.Context, id string) {
	c.JSON(200, &gen.OrderInfo{})
}
