package handlers

import (
	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
)

type DivindadeHandler struct {
	*GenericService
}

func NewDivindadeHandler() *DivindadeHandler {
	return &DivindadeHandler{
		GenericService: NewGenericService(database.DB),
	}
}

func (h *DivindadeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	divindades := rg.Group("/divindades")
	{
		divindades.GET("", h.GetAllDivindades)
		divindades.GET("/:id", h.GetDivindade)
		divindades.POST("", h.CreateDivindade)
		divindades.PUT("/:id", h.UpdateDivindade)
		divindades.DELETE("/:id", h.DeleteDivindade)
	}
}

func (h *DivindadeHandler) GetAllDivindades(c *gin.Context) {
	var divindades []models.Divindade
	h.GetAll(c, &divindades)
}

func (h *DivindadeHandler) GetDivindade(c *gin.Context) {
	var divindade models.Divindade
	h.GetByID(c, &divindade, "Divindade não encontrada")
}

func (h *DivindadeHandler) CreateDivindade(c *gin.Context) {
	var divindade models.Divindade
	h.Create(c, &divindade)
}

func (h *DivindadeHandler) UpdateDivindade(c *gin.Context) {
	var divindade models.Divindade
	h.Update(c, &divindade, "Divindade não encontrada")
}

func (h *DivindadeHandler) DeleteDivindade(c *gin.Context) {
	var divindade models.Divindade
	h.Delete(c, &divindade, "Divindade não encontrada")
}
