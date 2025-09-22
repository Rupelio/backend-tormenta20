package handlers

import (
	"net/http"
	"strconv"
	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
)

type ClasseHandler struct {
	*GenericService
}

func NewClasseHandler() *ClasseHandler {
	return &ClasseHandler{
		GenericService: NewGenericService(database.DB),
	}
}

func (h *ClasseHandler) RegisterRoutes(rg *gin.RouterGroup) {
	classes := rg.Group("/classes")
	{
		classes.GET("", h.GetAllClasses)
		classes.GET("/:id", h.GetClasse)
		classes.POST("", h.CreateClasse)
		classes.PUT("/:id", h.UpdateClasse)
		classes.DELETE("/:id", h.DeleteClasse)
		classes.PATCH("/:id/stats", h.UpdateClasseStats) // Nova rota para atualizar apenas PV e PM
	}
}

func (h *ClasseHandler) GetAllClasses(c *gin.Context) {
	var classes []models.Classe
	h.GetAll(c, &classes, "Habilidades")
}

func (h *ClasseHandler) GetClasse(c *gin.Context) {
	var classe models.Classe
	h.GetByID(c, &classe, "Classe não encontrada", "Habilidades")
}

func (h *ClasseHandler) CreateClasse(c *gin.Context) {
	var classe models.Classe
	h.Create(c, &classe, "Habilidades")
}

func (h *ClasseHandler) UpdateClasse(c *gin.Context) {
	var classe models.Classe
	h.Update(c, &classe, "Classe não encontrada", "Habilidades")
}

func (h *ClasseHandler) DeleteClasse(c *gin.Context) {
	var classe models.Classe
	h.Delete(c, &classe, "Classe não encontrada")
}

// UpdateClasseStats atualiza apenas os campos PV e PM de uma classe específica
func (h *ClasseHandler) UpdateClasseStats(c *gin.Context) {
	// Obter ID da URL
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Estrutura para receber apenas os campos de stats
	type ClasseStatsUpdate struct {
		PVPorNivel int `json:"pvpornivel" binding:"required,min=1,max=50"`
		PMPorNivel int `json:"pmpornivel" binding:"min=0,max=50"`
	}

	var statsUpdate ClasseStatsUpdate
	if err := c.ShouldBindJSON(&statsUpdate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar se a classe existe
	var classe models.Classe
	if err := database.DB.First(&classe, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Classe não encontrada"})
		return
	}

	// Atualizar apenas os campos de stats
	updateData := map[string]interface{}{
		"pv_por_nivel": statsUpdate.PVPorNivel,
		"pm_por_nivel": statsUpdate.PMPorNivel,
	}

	if err := database.DB.Model(&classe).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar stats da classe"})
		return
	}

	// Buscar a classe atualizada com as associações
	if err := database.DB.Preload("Habilidades").First(&classe, uint(id)).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar classe atualizada"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stats da classe atualizados com sucesso",
		"classe":  classe,
	})
}
