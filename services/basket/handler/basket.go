package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/puzzles/services/basket/dal"
	"github.com/puzzles/services/basket/gen"
)

type BasketHandler struct {
	db *dal.RedisDal
}

func NewBasketHandler(db *dal.RedisDal) *BasketHandler {
	return &BasketHandler{
		db: db,
	}
}
func (h *BasketHandler) CreateNewBasket(c *gin.Context) {
	var items []gen.Puzzle
	err := c.BindJSON(&items)
	if err != nil {
		c.JSON(400, gen.Error{Code: 400, Message: err.Error()})
	}

	id := uuid.New().String()
	basket := &gen.Basket{
		Id:    id,
		Items: items,
	}
	genErr := h.db.SetBasket(c, id, basket)
	if genErr != nil {
		c.JSON(int(genErr.Code), genErr)
	}

	c.JSON(201, basket)
}

func (h *BasketHandler) RemovePuzzleFromBasket(c *gin.Context, basketId string, puzzleId string) {
	c.Status(204)
}

func (h *BasketHandler) DeleteBasket(c *gin.Context, id string) {
	c.Status(204)
}

func (h *BasketHandler) GetBasket(c *gin.Context, id string) {
	basket, genErr := h.db.GetBasket(c, id)
	if genErr != nil {
		c.JSON(int(genErr.Code), genErr)
	}

	c.JSON(201, *basket)
}

func (h *BasketHandler) AddItemToBasket(c *gin.Context, id string) {
	c.JSON(200, gen.Basket{})
}
