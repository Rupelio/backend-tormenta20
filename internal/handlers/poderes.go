package handlers

import (
	"net/http"
	"strconv"

	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
)

type PoderHandler struct {
	*GenericService
}

func NewPoderHandler() *PoderHandler {
	return &PoderHandler{
		GenericService: NewGenericService(database.DB),
	}
}

func (h *PoderHandler) RegisterRoutes(rg *gin.RouterGroup) {
	poderes := rg.Group("/poderes")
	{
		poderes.GET("", h.GetAllPoderes)
		poderes.GET("/origem/:origem_id", h.GetPoderesPorOrigem)
		poderes.GET("/tipo/:tipo", h.GetPoderesPorTipo)
	}

	habilidades := rg.Group("/habilidades-especiais")
	{
		habilidades.GET("/raca/:raca_id", h.GetHabilidadesEspeciaisRaca)
	}
}

func (h *PoderHandler) GetAllPoderes(c *gin.Context) {
	var poderes []models.Poder
	h.GetAll(c, &poderes)
}

func (h *PoderHandler) GetPoderesPorOrigem(c *gin.Context) {
	origemIDStr := c.Param("origem_id")
	origemID, err := strconv.Atoi(origemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da origem inválido"})
		return
	}

	var poderes []models.Poder

	// Query para buscar poderes associados à origem
	err = database.DB.Table("poderes").
		Joins("INNER JOIN origem_poderes ON poderes.id = origem_poderes.poder_id").
		Where("origem_poderes.origem_id = ?", origemID).
		Find(&poderes).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar poderes da origem"})
		return
	}

	c.JSON(http.StatusOK, poderes)
}

func (h *PoderHandler) GetPoderesPorTipo(c *gin.Context) {
	tipo := c.Param("tipo")
	if tipo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tipo do poder não especificado"})
		return
	}

	var poderes []models.Poder
	err := database.DB.Where("tipo = ?", tipo).Find(&poderes).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar poderes por tipo"})
		return
	}

	c.JSON(http.StatusOK, poderes)
}

func (h *PoderHandler) GetHabilidadesEspeciaisRaca(c *gin.Context) {
	racaIDStr := c.Param("raca_id")
	racaID, err := strconv.Atoi(racaIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID da raça inválido"})
		return
	}

	var habilidades []models.RacaHabilidadeEspecial
	err = database.DB.Where("raca_id = ?", racaID).Find(&habilidades).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar habilidades especiais da raça"})
		return
	}

	c.JSON(http.StatusOK, habilidades)
}
