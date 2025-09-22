package handlers

import (
	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
)

type RacaHandler struct {
	*GenericService
}

func NewRacaHandler() *RacaHandler {
	return &RacaHandler{
		GenericService: NewGenericService(database.DB),
	}
}

func (h *RacaHandler) RegisterRoutes(rg *gin.RouterGroup) {
	racas := rg.Group("/racas")
	{
		racas.GET("", h.GetAllRacas)
		racas.GET("/:id", h.GetRaca)
		racas.POST("", h.CreateRaca)
		racas.PUT("/:id", h.UpdateRaca)
		racas.DELETE("/:id", h.DeleteRaca)
	}
}

func (h *RacaHandler) GetAllRacas(c *gin.Context) {
	var racas []models.Raca
	h.GetAll(c, &racas, "Habilidades")
}

func (h *RacaHandler) GetRaca(c *gin.Context) {
	var raca models.Raca
	h.GetByID(c, &raca, "Raça não encontrada", "Habilidades")
}

func (h *RacaHandler) CreateRaca(c *gin.Context) {
	var raca models.Raca
	h.Create(c, &raca, "Habilidades")
}

func (h *RacaHandler) UpdateRaca(c *gin.Context) {
	var raca models.Raca
	h.Update(c, &raca, "Raça não encontrada", "Habilidades")
}

func (h *RacaHandler) DeleteRaca(c *gin.Context) {
	var raca models.Raca
	h.Delete(c, &raca, "Raça não encontrada")
}
