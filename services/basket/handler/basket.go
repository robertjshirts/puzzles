package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/puzzles/services/basket/gen"
)

type BasketHandler struct {
}

func NewBasketHandler() *BasketHandler {
	return &BasketHandler{}
}

func (h *BasketHandler) CreateNewBasket(c *gin.Context) {
	c.JSON(200, gen.Basket{})
}

func (h *BasketHandler) RemovePuzzleFromBasket(c *gin.Context, basketId string, puzzleId string) {
	c.Status(204)
}

func (h *BasketHandler) DeleteBasket(c *gin.Context, id string) {
	c.Status(204)
}

func (h *BasketHandler) GetBasket(c *gin.Context, id string) {
	c.JSON(200, gen.Basket{})
}

func (h *BasketHandler) AddItemToBasket(c *gin.Context, id string) {
	c.JSON(200, gen.Basket{})
}
