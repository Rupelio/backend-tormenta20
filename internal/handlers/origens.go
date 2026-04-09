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
	if err := database.DB.Preload("Pericias").Preload("Itens").Preload("Habilidades").Find(&origens).Error; err != nil {
		h.Response.InternalError(c, "Erro ao buscar origens")
		return
	}
	c.JSON(200, origens)
}

func (h *OrigemHandler) GetOrigem(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}
	var origem models.Origem
	if err := database.DB.Preload("Pericias").Preload("Itens").Preload("Habilidades").First(&origem, id).Error; err != nil {
		h.Response.NotFound(c, "Origem não encontrada")
		return
	}
	h.Response.Success(c, origem)
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
