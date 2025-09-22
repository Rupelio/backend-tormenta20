package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"tormenta20-builder/internal/models"
	"gorm.io/gorm"
)

type PericiasHandler struct {
	db *gorm.DB
}

func NewPericiasHandler(db *gorm.DB) *PericiasHandler {
	return &PericiasHandler{db: db}
}

// GET /pericias - Lista todas as perícias
func (h *PericiasHandler) GetPericias(c *gin.Context) {
	var pericias []models.Pericia

	if err := h.db.Find(&pericias).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perícias"})
		return
	}

	c.JSON(http.StatusOK, pericias)
}

// GET /pericias/:id - Busca uma perícia por ID
func (h *PericiasHandler) GetPericia(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var pericia models.Pericia
	if err := h.db.First(&pericia, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Perícia não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perícia"})
		return
	}

	c.JSON(http.StatusOK, pericia)
}

// GET /classes/:id/pericias - Lista perícias disponíveis para uma classe
func (h *PericiasHandler) GetPericiasClasse(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var classe models.Classe
	if err := h.db.Preload("PericiasDisponiveis").Preload("PericiasAutomaticas").First(&classe, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Classe não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perícias da classe"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"classe": classe.Nome,
		"pericias_quantidade": classe.PericiasQuantidade,
		"pericias_disponiveis": classe.PericiasDisponiveis,
		"pericias_automaticas": classe.PericiasAutomaticas,
	})
}

// GET /racas/:id/pericias - Lista perícias de uma raça
func (h *PericiasHandler) GetPericiasRaca(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var raca models.Raca
	if err := h.db.Preload("Pericias").First(&raca, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Raça não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perícias da raça"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"raca": raca.Nome,
		"pericias": raca.Pericias,
	})
}

// GET /origens/:id/pericias - Lista perícias de uma origem
func (h *PericiasHandler) GetPericiasOrigem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var origem models.Origem
	if err := h.db.Preload("Pericias").First(&origem, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Origem não encontrada"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perícias da origem"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"origem": origem.Nome,
		"pericias": origem.Pericias,
	})
}

// POST /personagens/:id/pericias - Atualiza perícias do personagem
func (h *PericiasHandler) UpdatePericiasPersonagem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	var request struct {
		PericiasIDs []uint `json:"pericias_ids"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dados inválidos"})
		return
	}

	// Buscar personagem com suas relações
	var personagem models.Personagem
	if err := h.db.Preload("Classe").Preload("Raca").Preload("Origem").First(&personagem, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar personagem"})
		return
	}

	// Buscar perícias disponíveis da classe
	var classe models.Classe
	if err := h.db.Preload("PericiasDisponiveis").Preload("PericiasAutomaticas").First(&classe, personagem.ClasseID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perícias da classe"})
		return
	}

	// Buscar perícias automáticas da raça e origem
	var raca models.Raca
	var origem models.Origem
	var periciasAutomaticas []models.Pericia

	if err := h.db.Preload("Pericias").First(&raca, personagem.RacaID).Error; err == nil {
		periciasAutomaticas = append(periciasAutomaticas, raca.Pericias...)
	}

	if err := h.db.Preload("Pericias").First(&origem, personagem.OrigemID).Error; err == nil {
		periciasAutomaticas = append(periciasAutomaticas, origem.Pericias...)
	}

	// Adicionar perícias automáticas da classe
	periciasAutomaticas = append(periciasAutomaticas, classe.PericiasAutomaticas...)

	// Validar que as perícias escolhidas estão disponíveis para a classe
	for _, periciaID := range request.PericiasIDs {
		disponivel := false
		for _, pericia := range classe.PericiasDisponiveis {
			if pericia.ID == periciaID {
				disponivel = true
				break
			}
		}
		if !disponivel {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("Perícia ID %d não está disponível para a classe %s", periciaID, classe.Nome),
			})
			return
		}

		// Verificar se não é uma perícia automática
		for _, automatica := range periciasAutomaticas {
			if automatica.ID == periciaID {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": fmt.Sprintf("Perícia %s já é automática da raça/origem", automatica.Nome),
				})
				return
			}
		}
	}

	// Validar quantidade de perícias
	if len(request.PericiasIDs) > classe.PericiasQuantidade {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Máximo de %d perícias permitidas para %s", classe.PericiasQuantidade, classe.Nome),
		})
		return
	}

	// Buscar perícias selecionadas
	var pericias []models.Pericia
	if err := h.db.Find(&pericias, request.PericiasIDs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perícias"})
		return
	}

	// Verificar se todas as perícias foram encontradas
	if len(pericias) != len(request.PericiasIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Uma ou mais perícias não foram encontradas"})
		return
	}

	// Substituir perícias do personagem
	if err := h.db.Model(&personagem).Association("Pericias").Replace(pericias); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar perícias"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Perícias atualizadas com sucesso",
		"pericias_escolhidas": len(pericias),
		"pericias_automaticas": len(periciasAutomaticas),
	})
}
