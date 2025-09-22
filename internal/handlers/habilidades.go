package handlers

import (
	"net/http"
	"strconv"
	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
)

type HabilidadeHandler struct {
	*GenericService
}

func NewHabilidadeHandler() *HabilidadeHandler {
	return &HabilidadeHandler{
		GenericService: NewGenericService(database.DB),
	}
}

func (h *HabilidadeHandler) RegisterRoutes(rg *gin.RouterGroup) {
	habilidades := rg.Group("/habilidades")
	{
		// Rotas para habilidades de raça
		habilidades.GET("/raca/:id", h.GetHabilidadesRaca)

		// Rotas para habilidades de classe
		habilidades.GET("/classe/:id", h.GetHabilidadesClasse)
		habilidades.GET("/classe/:id/nivel/:nivel", h.GetHabilidadesClassePorNivel)

		// Rotas para habilidades de origem
		habilidades.GET("/origem/:id", h.GetHabilidadesOrigem)

		// Rotas para habilidades de divindade
		habilidades.GET("/divindade/:id", h.GetHabilidadesDivindade)
		habilidades.GET("/divindade/:id/nivel/:nivel", h.GetHabilidadesDivindadePorNivel)
	}
}

// GetHabilidadesRaca retorna todas as habilidades de uma raça específica
func (h *HabilidadeHandler) GetHabilidadesRaca(c *gin.Context) {
	id := c.Param("id")
	racaID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var habilidades []models.HabilidadeRaca
	if err := database.DB.Where("raca_id = ?", racaID).Find(&habilidades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar habilidades da raça"})
		return
	}

	c.JSON(http.StatusOK, habilidades)
}

// GetHabilidadesClasse retorna todas as habilidades de uma classe específica
func (h *HabilidadeHandler) GetHabilidadesClasse(c *gin.Context) {
	id := c.Param("id")
	classeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var habilidades []models.HabilidadeClasse
	if err := database.DB.Where("classe_id = ?", classeID).Order("nivel").Find(&habilidades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar habilidades da classe"})
		return
	}

	c.JSON(http.StatusOK, habilidades)
}

// GetHabilidadesClassePorNivel retorna habilidades de uma classe até um nível específico
func (h *HabilidadeHandler) GetHabilidadesClassePorNivel(c *gin.Context) {
	id := c.Param("id")
	nivel := c.Param("nivel")

	classeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	nivelInt, err := strconv.Atoi(nivel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nível inválido"})
		return
	}

	var habilidades []models.HabilidadeClasse
	if err := database.DB.Where("classe_id = ? AND nivel <= ?", classeID, nivelInt).Order("nivel").Find(&habilidades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar habilidades da classe"})
		return
	}

	c.JSON(http.StatusOK, habilidades)
}

// GetHabilidadesOrigem retorna todas as habilidades de uma origem específica
func (h *HabilidadeHandler) GetHabilidadesOrigem(c *gin.Context) {
	id := c.Param("id")
	origemID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var habilidades []models.HabilidadeOrigem
	if err := database.DB.Where("origem_id = ?", origemID).Find(&habilidades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar habilidades da origem"})
		return
	}

	c.JSON(http.StatusOK, habilidades)
}

// GetHabilidadesDivindade retorna todas as habilidades de uma divindade específica
func (h *HabilidadeHandler) GetHabilidadesDivindade(c *gin.Context) {
	id := c.Param("id")
	divindadeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var habilidades []models.HabilidadeDivindade
	if err := database.DB.Where("divindade_id = ?", divindadeID).Order("nivel").Find(&habilidades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar habilidades da divindade"})
		return
	}

	c.JSON(http.StatusOK, habilidades)
}

// GetHabilidadesDivindadePorNivel retorna habilidades de uma divindade até um nível específico
func (h *HabilidadeHandler) GetHabilidadesDivindadePorNivel(c *gin.Context) {
	id := c.Param("id")
	nivel := c.Param("nivel")

	divindadeID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	nivelInt, err := strconv.Atoi(nivel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nível inválido"})
		return
	}

	var habilidades []models.HabilidadeDivindade
	if err := database.DB.Where("divindade_id = ? AND nivel <= ?", divindadeID, nivelInt).Order("nivel").Find(&habilidades).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar habilidades da divindade"})
		return
	}

	c.JSON(http.StatusOK, habilidades)
}
