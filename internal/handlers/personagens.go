package handlers

import (
	"net/http"

	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/middleware"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PersonagemHandler struct {
	*GenericService
}

func NewPersonagemHandler() *PersonagemHandler {
	return &PersonagemHandler{
		GenericService: NewGenericService(database.DB),
	}
}

func (h *PersonagemHandler) RegisterRoutes(rg *gin.RouterGroup) {
	personagens := rg.Group("/personagens")
	{
		personagens.GET("", h.GetAllPersonagens)
		personagens.GET("/:id", h.GetPersonagem)
		personagens.POST("", h.CreatePersonagem)
		personagens.PUT("/:id", h.UpdatePersonagem)
		personagens.DELETE("/:id", h.DeletePersonagem)
		personagens.POST("/calculate", h.CalculateStats)
		personagens.GET("/:id/export-pdf", h.ExportToPDF)
		personagens.GET("/:id/test", h.TestRoute) // Rota de teste
	}
}

func (h *PersonagemHandler) GetAllPersonagens(c *gin.Context) {
	var personagens []models.Personagem

	// Obtém identificação do usuário
	sessionID, userIP := middleware.GetUserIdentification(c)

	// Constrói query para filtrar personagens do usuário
	query := database.DB.Preload("Raca").Preload("Classe").Preload("Origem")

	if sessionID != "" {
		// Prioriza session ID se disponível
		query = query.Where("user_session_id = ?", sessionID)
	} else if userIP != "" {
		// Fallback para IP se não há session
		query = query.Where("user_ip = ?", userIP)
	} else {
		// Se não há identificação, retorna vazio (não deve acontecer com middleware)
		c.JSON(http.StatusOK, []models.Personagem{})
		return
	}

	if err := query.Find(&personagens).Error; err != nil {
		h.Response.InternalError(c, "Erro ao buscar personagens: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, personagens)
}

func (h *PersonagemHandler) GetPersonagem(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	var personagem models.Personagem
	sessionID, userIP := middleware.GetUserIdentification(c)

	// Constrói query para buscar personagem do usuário
	query := database.DB.Preload("Raca").Preload("Classe").Preload("Origem")

	if sessionID != "" {
		query = query.Where("id = ? AND user_session_id = ?", id, sessionID)
	} else if userIP != "" {
		query = query.Where("id = ? AND user_ip = ?", id, userIP)
	} else {
		h.Response.NotFound(c, "Personagem não encontrado")
		return
	}

	if err := query.First(&personagem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	h.Response.Success(c, personagem)
}

func (h *PersonagemHandler) CreatePersonagem(c *gin.Context) {
	var personagem models.Personagem

	// Bind dos dados do personagem
	if err := c.ShouldBindJSON(&personagem); err != nil {
		h.Response.BadRequest(c, err.Error())
		return
	}

	// Validar estrutura
	if err := validate.Struct(personagem); err != nil {
		h.Response.BadRequest(c, err.Error())
		return
	}

	// Obtém identificação do usuário e adiciona ao personagem
	sessionID, userIP := middleware.GetUserIdentification(c)

	if sessionID != "" {
		personagem.UserSessionID = &sessionID
		personagem.CreatedByType = "session"
	}

	if userIP != "" {
		personagem.UserIP = &userIP
		// Se tem session E IP, define como hybrid
		if sessionID != "" {
			personagem.CreatedByType = "hybrid"
		} else {
			personagem.CreatedByType = "ip"
		}
	}

	// Garante que escolhas_raca seja um JSON válido
	if personagem.EscolhasRaca == "" {
		personagem.EscolhasRaca = "{}"
	}

	// Criar no banco
	if err := database.DB.Create(&personagem).Error; err != nil {
		h.Response.InternalError(c, "Erro ao criar personagem")
		return
	}

	// Recarregar com relações
	if err := database.DB.Preload("Raca").Preload("Classe").Preload("Origem").First(&personagem, personagem.ID).Error; err != nil {
		h.Response.InternalError(c, "Erro ao carregar personagem criado")
		return
	}

	h.Response.Created(c, personagem)
}

func (h *PersonagemHandler) UpdatePersonagem(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	var personagem models.Personagem
	sessionID, userIP := middleware.GetUserIdentification(c)

	// Primeiro, verifica se o personagem existe e pertence ao usuário
	query := database.DB
	if sessionID != "" {
		query = query.Where("id = ? AND user_session_id = ?", id, sessionID)
	} else if userIP != "" {
		query = query.Where("id = ? AND user_ip = ?", id, userIP)
	} else {
		h.Response.NotFound(c, "Personagem não encontrado")
		return
	}

	if err := query.First(&personagem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	// Bind dos novos dados (preservando a identificação do usuário)
	originalSessionID := personagem.UserSessionID
	originalUserIP := personagem.UserIP
	originalCreatedByType := personagem.CreatedByType

	if err := c.ShouldBindJSON(&personagem); err != nil {
		h.Response.BadRequest(c, err.Error())
		return
	}

	// Restaura identificação original (não deve ser alterada)
	personagem.UserSessionID = originalSessionID
	personagem.UserIP = originalUserIP
	personagem.CreatedByType = originalCreatedByType

	// Validar estrutura
	if err := validate.Struct(personagem); err != nil {
		h.Response.BadRequest(c, err.Error())
		return
	}

	// Salvar alterações
	if err := database.DB.Save(&personagem).Error; err != nil {
		h.Response.InternalError(c, "Erro ao atualizar personagem")
		return
	}

	// Recarregar com relações
	if err := database.DB.Preload("Raca").Preload("Classe").Preload("Origem").First(&personagem, personagem.ID).Error; err != nil {
		h.Response.InternalError(c, "Erro ao carregar personagem atualizado")
		return
	}

	h.Response.Success(c, personagem)
}

func (h *PersonagemHandler) DeletePersonagem(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	var personagem models.Personagem
	sessionID, userIP := middleware.GetUserIdentification(c)

	// Verifica se o personagem existe e pertence ao usuário
	query := database.DB
	if sessionID != "" {
		query = query.Where("id = ? AND user_session_id = ?", id, sessionID)
	} else if userIP != "" {
		query = query.Where("id = ? AND user_ip = ?", id, userIP)
	} else {
		h.Response.NotFound(c, "Personagem não encontrado")
		return
	}

	if err := query.First(&personagem).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	if err := database.DB.Delete(&personagem).Error; err != nil {
		h.Response.InternalError(c, "Erro ao deletar personagem")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// CalculateStats calcula as estatísticas de um personagem baseado em seus atributos, raça e classe
func (h *PersonagemHandler) CalculateStats(c *gin.Context) {
	var personagemData struct {
		Nivel        int `json:"nivel"`
		Forca        int `json:"forca"`
		Destreza     int `json:"destreza"`
		Constituicao int `json:"constituicao"`
		Inteligencia int `json:"inteligencia"`
		Sabedoria    int `json:"sabedoria"`
		Carisma      int `json:"carisma"`
		RacaID       int `json:"raca_id"`
		ClasseID     int `json:"classe_id"`
		OrigemID     int `json:"origem_id"`
	}

	if err := c.ShouldBindJSON(&personagemData); err != nil {
		h.Response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	// Buscar dados da raça e classe para aplicar modificadores
	var raca models.Raca
	var classe models.Classe

	if err := database.DB.First(&raca, personagemData.RacaID).Error; err != nil {
		h.Response.NotFound(c, "Raça não encontrada")
		return
	}

	if err := database.DB.First(&classe, personagemData.ClasseID).Error; err != nil {
		h.Response.NotFound(c, "Classe não encontrada")
		return
	}

	// Calcular modificadores dos atributos (Tormenta20: modificador = valor do atributo)
	calcularModificador := func(valor int) int {
		return valor
	}

	modCon := calcularModificador(personagemData.Constituicao)
	modInt := calcularModificador(personagemData.Inteligencia)
	modDes := calcularModificador(personagemData.Destreza)

	// Calcular PV (Pontos de Vida)
	// PV = (PV por nível da classe + mod CON) × nível + bônus racial
	pvTotal := (classe.PVPorNivel + modCon) * personagemData.Nivel

	// Calcular PM (Pontos de Mana)
	// PM = (PM por nível da classe + mod INT) × nível + bônus racial
	pmTotal := (classe.PMPorNivel + modInt) * personagemData.Nivel

	// Calcular Defesa
	// Defesa = 10 + mod DES + bônus de armadura (assumindo 0) + bônus racial
	defesa := 10 + modDes

	// Aplicar modificadores raciais se existirem
	// (Aqui você pode adicionar lógica específica para cada raça)
	switch raca.Nome {
	case "Humano":
		// Humanos ganham +1 em qualquer atributo (já considerado na criação)
	case "Elfo":
		// Elfos podem ter bônus específicos
		defesa += 1 // Exemplo: +1 na defesa
	case "Anão":
		// Anões podem ter bônus de PV
		pvTotal += personagemData.Nivel * 2
	}

	// Aplicar modificadores de classe se necessário
	switch classe.Nome {
	case "Guerreiro":
		pvTotal += personagemData.Nivel // Guerreiros ganham PV extra
	case "Mago":
		pmTotal += personagemData.Nivel // Magos ganham PM extra
	}

	stats := gin.H{
		"pv_total": pvTotal,
		"pm_total": pmTotal,
		"defesa":   defesa,
	}

	c.JSON(http.StatusOK, stats)
}

// ExportToPDF exporta a ficha do personagem em formato PDF
func (h *PersonagemHandler) ExportToPDF(c *gin.Context) {
	id := c.Param("id")

	var personagem models.Personagem
	if err := database.DB.Preload("Raca").Preload("Classe").Preload("Origem").First(&personagem, id).Error; err != nil {
		h.Response.NotFound(c, "Personagem não encontrado")
		return
	}

	// Resposta temporária simplificada para testar se a rota funciona
	c.JSON(http.StatusOK, gin.H{
		"message":    "PDF export funcionando",
		"personagem": personagem.Nome,
		"id":         id,
	})
}

// TestRoute rota de teste para verificar se funciona
func (h *PersonagemHandler) TestRoute(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"message": "Rota de teste funcionando",
		"id":      id,
	})
}
