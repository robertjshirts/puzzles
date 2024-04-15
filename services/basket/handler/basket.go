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
		return
	}

	id := uuid.New().String()
	basket := &gen.Basket{
		Id:    id,
		Items: items,
	}
	genErr := h.db.SetBasket(c, id, basket)
	if genErr != nil {
		c.JSON(int(genErr.Code), genErr)
		return
	}

	c.JSON(201, basket)
}

func (h *BasketHandler) RemovePuzzleFromBasket(c *gin.Context, basketId string, puzzleId string) {
	basket, err := h.db.GetBasket(c, basketId)
	if err != nil {
		c.JSON(int(err.Code), err)
		return
	}

	for i, item := range basket.Items {
		if item.Id == puzzleId {
			basket.Items = append(basket.Items[:i], basket.Items[i+1:]...)
			break
		}
	}

	genErr := h.db.SetBasket(c, basketId, basket)
	if genErr != nil {
		c.JSON(int(genErr.Code), genErr)
		return
	}

	c.Status(204)
}

func (h *BasketHandler) DeleteBasket(c *gin.Context, id string) {
	err := h.db.DeleteBasket(c, id)
	if err != nil {
		c.JSON(int(err.Code), err)
		return
	}

	c.Status(204)
}

func (h *BasketHandler) GetBasket(c *gin.Context, id string) {
	basket, genErr := h.db.GetBasket(c, id)
	if genErr != nil {
		c.JSON(int(genErr.Code), genErr)
		return
	}

	c.JSON(200, basket)
}

func (h *BasketHandler) AddItemToBasket(c *gin.Context, id string) {
	var items []gen.Puzzle
	err := c.BindJSON(&items)
	if err != nil {
		c.JSON(400, gen.Error{Code: 400, Message: err.Error()})
		return
	}

	basket, genErr := h.db.GetBasket(c, id)
	if genErr != nil {
		c.JSON(int(genErr.Code), genErr)
		return
	}
	basket.Items = append(basket.Items, items...)

	genErr = h.db.SetBasket(c, id, basket)
	if genErr != nil {
		c.JSON(int(genErr.Code), genErr)
		return
	}

	c.JSON(201, basket)
}
