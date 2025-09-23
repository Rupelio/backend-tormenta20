package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
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
		"classe":               classe.Nome,
		"pericias_quantidade":  classe.PericiasQuantidade,
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
		"raca":     raca.Nome,
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
		"origem":   origem.Nome,
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

	// Separar perícias de classe das perícias de raça
	var periciasClasse []uint
	var periciasRaca []uint

	// Buscar perícias automáticas da raça e origem para identificar quais são de raça
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

	// Separar perícias enviadas entre classe e raça
	for _, periciaID := range request.PericiasIDs {
		// Verificar se é uma perícia disponível para a classe
		disponivelParaClasse := false
		for _, pericia := range classe.PericiasDisponiveis {
			if pericia.ID == periciaID {
				disponivelParaClasse = true
				break
			}
		}

		// Verificar se é uma perícia automática (que não deve ser validada)
		ehAutomatica := false
		for _, automatica := range periciasAutomaticas {
			if automatica.ID == periciaID {
				ehAutomatica = true
				break
			}
		}

		if ehAutomatica {
			// Perícias automáticas não são salvas como escolhas do usuário
			continue
		} else if disponivelParaClasse {
			// Perícia disponível para a classe - adicionar às perícias de classe
			periciasClasse = append(periciasClasse, periciaID)
		} else {
			// Perícia não disponível para classe - assumir que é de raça/origem
			periciasRaca = append(periciasRaca, periciaID)
		}
	}

	// Validar quantidade de perícias DE CLASSE apenas
	if len(periciasClasse) > classe.PericiasQuantidade {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("Máximo de %d perícias de classe permitidas para %s (recebidas %d)", classe.PericiasQuantidade, classe.Nome, len(periciasClasse)),
		})
		return
	}

	// Buscar todas as perícias selecionadas para verificar se existem
	var todasPericias []uint
	todasPericias = append(todasPericias, periciasClasse...)
	todasPericias = append(todasPericias, periciasRaca...)

	var pericias []models.Pericia
	if len(todasPericias) > 0 {
		if err := h.db.Find(&pericias, todasPericias).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar perícias"})
			return
		}

		// Verificar se todas as perícias foram encontradas
		if len(pericias) != len(todasPericias) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Uma ou mais perícias não foram encontradas"})
			return
		}
	}

	// Limpar perícias escolhidas existentes (tanto de classe quanto de raça)
	if err := h.db.Where("personagem_id = ? AND fonte IN (?)", id, []string{"classe", "raca"}).Delete(&models.PersonagemPericia{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao limpar perícias existentes"})
		return
	}

	// Inserir perícias de classe
	for _, periciaID := range periciasClasse {
		personagemPericia := models.PersonagemPericia{
			PersonagemID: uint(id),
			PericiaID:    periciaID,
			Fonte:        "classe",
		}
		if err := h.db.Create(&personagemPericia).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar perícias de classe"})
			return
		}
	}

	// Inserir perícias de raça
	for _, periciaID := range periciasRaca {
		personagemPericia := models.PersonagemPericia{
			PersonagemID: uint(id),
			PericiaID:    periciaID,
			Fonte:        "raca",
		}
		if err := h.db.Create(&personagemPericia).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar perícias de raça"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":              "Perícias atualizadas com sucesso",
		"pericias_classe":      len(periciasClasse),
		"pericias_raca":        len(periciasRaca),
		"pericias_automaticas": len(periciasAutomaticas),
	})
}

// GET /personagens/:id/pericias - Busca perícias selecionadas do personagem
func (h *PericiasHandler) GetPericiasPersonagem(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID inválido"})
		return
	}

	// Buscar personagem com suas perícias
	var personagem models.Personagem
	if err := h.db.Preload("Pericias").First(&personagem, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Personagem não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar personagem"})
		return
	}

	// Extrair IDs das perícias
	var periciasIds []uint
	for _, pericia := range personagem.Pericias {
		periciasIds = append(periciasIds, pericia.ID)
	}

	c.JSON(http.StatusOK, gin.H{
		"personagem_id": id,
		"pericias_ids":  periciasIds,
	})
}
