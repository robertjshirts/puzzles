package api

import (
	"github.com/gin-gonic/gin"

	"github.com/puzzles/services/catalog/dal"
	"github.com/puzzles/services/catalog/gen"
)

type CatalogHandler struct {
	db *dal.MongoDAL
}

func NewCatalogHandler(db *dal.MongoDAL) *CatalogHandler {
	return &CatalogHandler{db: db}
}

func (h *CatalogHandler) GetPuzzles(c *gin.Context) {
	puzzles, err := h.db.GetAllPuzzles(c)
	if err != nil {
		c.JSON(500, gen.Error{Code: 500, Message: err.Error()})
	}
	c.JSON(200, puzzles)
}

func (h *CatalogHandler) AddPuzzle(c *gin.Context) {
	var newPuzzle gen.NewPuzzle
	if err := c.BindJSON(&newPuzzle); err != nil {
		c.JSON(400, gen.Error{Code: 400, Message: err.Error()})
	}

	puzzle, err := h.db.AddPuzzle(c, newPuzzle)
	if err != nil {
		c.JSON(500, gen.Error{Code: 500, Message: err.Error()})
	}
	c.JSON(201, puzzle)
}

func (h *CatalogHandler) DeletePuzzle(c *gin.Context, id string) {
	if err := h.db.DeletePuzzle(c, id); err != nil {
		c.JSON(500, gen.Error{Code: 500, Message: err.Error()})
	}
	c.JSON(204, nil)
}

func (h *CatalogHandler) GetPuzzle(c *gin.Context, id string) {
	puzzle, err := h.db.GetPuzzle(c, id)
	if err != nil {
		c.JSON(500, gen.Error{Code: 500, Message: err.Error()})
	}
	c.JSON(200, puzzle)
}

func (*CatalogHandler) UpdatePuzzle(c *gin.Context, id string) {
	c.JSON(200, "UpdatePuzzle not implemented")
}
