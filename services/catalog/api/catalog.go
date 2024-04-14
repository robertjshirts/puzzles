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
		c.JSON(int(err.Code), err)
	}
	c.JSON(200, puzzles)
}

func (h *CatalogHandler) AddPuzzle(c *gin.Context) {
	var newPuzzle gen.NewPuzzle
	err := c.BindJSON(&newPuzzle)
	if err != nil {
		c.JSON(400, gen.Error{Code: 400, Message: err.Error()})
	}

	puzzle, genErr := h.db.AddPuzzle(c, newPuzzle)
	if genErr != nil {
		c.JSON(int(genErr.Code), err)
	}
	c.JSON(201, puzzle)
}

func (h *CatalogHandler) DeletePuzzle(c *gin.Context, id string) {
	if err := h.db.DeletePuzzle(c, id); err != nil {
		c.JSON(500, err)
	}
	c.JSON(204, nil)
}

func (h *CatalogHandler) GetPuzzle(c *gin.Context, id string) {
	puzzle, err := h.db.GetPuzzle(c, id)
	if err != nil {
		c.JSON(int(err.Code), err)
	}
	c.JSON(200, puzzle)
}

func (h *CatalogHandler) UpdatePuzzle(c *gin.Context, id string) {
	var updates gen.PuzzleUpdate
	if err := c.BindJSON(&updates); err != nil {
		c.JSON(400, gen.Error{Code: 400, Message: err.Error()})
	}
	err := h.db.UpdatePuzzle(c, id, updates)
	if err != nil {
		c.JSON(int(err.Code), err)
	}

	c.JSON(200, nil)
}
