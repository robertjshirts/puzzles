package handler

import (
	"github.com/puzzles/services/orders/gen"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) CreateNewOrder(c *gin.Context) {
	c.JSON(201, &gen.OrderInfo{})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context, id string) {
	c.Status(204)
}

func (h *OrderHandler) GetOrder(c *gin.Context, id string) {
	c.JSON(200, &gen.OrderInfo{})
}
