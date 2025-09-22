package handlers

import (
	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
)

type OrigemHandler struct {
	*GenericService
}

func NewOrigemHandler() *OrigemHandler {
	return &OrigemHandler{
		GenericService: NewGenericService(database.DB),
	}
}

func (h *OrigemHandler) RegisterRoutes(rg *gin.RouterGroup) {
	origens := rg.Group("/origens")
	{
		origens.GET("", h.GetAllOrigens)
		origens.GET("/:id", h.GetOrigem)
		origens.POST("", h.CreateOrigem)
		origens.PUT("/:id", h.UpdateOrigem)
		origens.DELETE("/:id", h.DeleteOrigem)
	}
}

func (h *OrigemHandler) GetAllOrigens(c *gin.Context) {
	var origens []models.Origem
	h.GetAll(c, &origens)
}

func (h *OrigemHandler) GetOrigem(c *gin.Context) {
	var origem models.Origem
	h.GetByID(c, &origem, "Origem não encontrada")
}

func (h *OrigemHandler) CreateOrigem(c *gin.Context) {
	var origem models.Origem
	h.Create(c, &origem)
}

func (h *OrigemHandler) UpdateOrigem(c *gin.Context) {
	var origem models.Origem
	h.Update(c, &origem, "Origem não encontrada")
}

func (h *OrigemHandler) DeleteOrigem(c *gin.Context) {
	var origem models.Origem
	h.Delete(c, &origem, "Origem não encontrada")
}
