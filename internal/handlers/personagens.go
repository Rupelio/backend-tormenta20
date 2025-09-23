package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/middleware"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PersonagemRequest representa os dados recebidos do frontend
type PersonagemRequest struct {
	// Dados do personagem
	Nome         string `json:"nome" validate:"required,min=2"`
	Nivel        int    `json:"nivel" validate:"min=1,max=20"`
	For          int    `json:"for" validate:"min=0,max=4"`
	Des          int    `json:"des" validate:"min=0,max=4"`
	Con          int    `json:"con" validate:"min=0,max=4"`
	Int          int    `json:"int" validate:"min=0,max=4"`
	Sab          int    `json:"sab" validate:"min=0,max=4"`
	Car          int    `json:"car" validate:"min=0,max=4"`
	RacaID       uint   `json:"raca_id"`
	ClasseID     uint   `json:"classe_id"`
	OrigemID     uint   `json:"origem_id"`
	DivindadeID  *uint  `json:"divindade_id"`
	EscolhasRaca string `json:"escolhas_raca"`

	// Dados complementares
	AtributosLivres      []string `json:"atributosLivres"`       // Array de atributos livres escolhidos
	PericiasSelecionadas []uint   `json:"pericias_selecionadas"` // IDs das perícias selecionadas
	PoderesClasse        []uint   `json:"poderes_classe"`        // IDs dos poderes de classe
	PoderesDivinos       []uint   `json:"poderes_divinos"`       // IDs dos poderes divinos
}

type PersonagemHandler struct {
	*GenericService
}

func NewPersonagemHandler() *PersonagemHandler {
	return &PersonagemHandler{
		GenericService: NewGenericService(database.DB),
	}
}

// findPersonagemByUser busca um personagem que pertence ao usuário (por sessão OU IP)
func (h *PersonagemHandler) findPersonagemByUser(c *gin.Context, personagemID int) (*models.Personagem, error) {
	sessionID, userIP := middleware.GetUserIdentification(c)

	var personagem models.Personagem
	query := database.DB

	if sessionID != "" && userIP != "" {
		// Se tem ambos, busca por qualquer um dos dois (para lidar com sessões que mudam)
		query = query.Where("id = ? AND (user_session_id = ? OR user_ip = ?)", personagemID, sessionID, userIP)
	} else if sessionID != "" {
		query = query.Where("id = ? AND user_session_id = ?", personagemID, sessionID)
	} else if userIP != "" {
		query = query.Where("id = ? AND user_ip = ?", personagemID, userIP)
	} else {
		return nil, gorm.ErrRecordNotFound
	}

	err := query.First(&personagem).Error
	return &personagem, err
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
		// Novos endpoints para poderes
		personagens.POST("/:id/poderes-divinos", h.SavePoderesDivinos)
		personagens.POST("/:id/poderes-classe", h.SavePoderesClasse)
		personagens.GET("/:id/poderes-divinos", h.GetPoderesDivinos)
		personagens.GET("/:id/poderes-classe", h.GetPoderesClasse)
		// Endpoint para escolhas raciais
		personagens.POST("/:id/escolhas-raca", h.SaveEscolhasRaca)
		personagens.GET("/:id/escolhas-raca", h.GetEscolhasRaca)
		// Endpoint de debug para ver TODOS os personagens (sem filtro de usuário)
		// personagens.GET("/debug/all", h.GetAllPersonagensDebug)
	}
}

func (h *PersonagemHandler) GetAllPersonagens(c *gin.Context) {
	var personagens []models.Personagem

	// Obtém identificação do usuário
	sessionID, userIP := middleware.GetUserIdentification(c)

	// Constrói query para filtrar personagens do usuário
	query := database.DB.Preload("Raca").Preload("Classe").Preload("Origem").Preload("Divindade")

	if sessionID != "" && userIP != "" {
		// Se tem ambos, busca por qualquer um dos dois (para lidar com sessões que mudam)
		query = query.Where("user_session_id = ? OR user_ip = ?", sessionID, userIP)
	} else if sessionID != "" {
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

	// Carregar perícias manualmente e calcular stats para cada personagem
	for i := range personagens {
		h.loadPersonagemPericias(&personagens[i])
		h.calculatePersonagemStats(&personagens[i])
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
	query := database.DB.Preload("Raca").Preload("Classe").Preload("Origem").Preload("Divindade")

	if sessionID != "" && userIP != "" {
		// Se tem ambos, busca por qualquer um dos dois (para lidar com sessões que mudam)
		query = query.Where("id = ? AND (user_session_id = ? OR user_ip = ?)", id, sessionID, userIP)
	} else if sessionID != "" {
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

	// Carregar perícias manualmente e calcular stats do personagem
	h.loadPersonagemPericias(&personagem)
	h.calculatePersonagemStats(&personagem)

	h.Response.Success(c, personagem)
}

func (h *PersonagemHandler) CreatePersonagem(c *gin.Context) {
	var req PersonagemRequest

	// Log do payload recebido para debug
	body, _ := c.GetRawData()
	fmt.Printf("DEBUG - CreatePersonagem - Raw body: %s\n", string(body))

	// Recriar o context para poder fazer bind novamente
	c.Request.Body = io.NopCloser(strings.NewReader(string(body)))

	// Bind dos dados do personagem
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Printf("DEBUG - CreatePersonagem - Erro no bind JSON: %v\n", err)
		h.Response.BadRequest(c, err.Error())
		return
	}

	fmt.Printf("DEBUG - CreatePersonagem - Dados parseados: %+v\n", req)

	// Criar personagem a partir dos dados da request
	personagem := models.Personagem{
		Nome:        req.Nome,
		Nivel:       req.Nivel,
		For:         req.For,
		Des:         req.Des,
		Con:         req.Con,
		Int:         req.Int,
		Sab:         req.Sab,
		Car:         req.Car,
		RacaID:      req.RacaID,
		ClasseID:    req.ClasseID,
		OrigemID:    req.OrigemID,
		DivindadeID: req.DivindadeID,
	}

	// Garantir que EscolhasRaca seja um JSON válido
	if req.EscolhasRaca == "" {
		personagem.EscolhasRaca = "{}"
	} else {
		personagem.EscolhasRaca = req.EscolhasRaca
	}

	// LIMPAR ID para evitar conflito de chave primária na criação
	personagem.ID = 0

	// LIMPAR relações many2many para evitar problemas no Create
	personagem.Pericias = nil

	// Processar atributos livres
	if len(req.AtributosLivres) > 0 {
		atributosLivresJSON, err := json.Marshal(req.AtributosLivres)
		if err != nil {
			h.Response.BadRequest(c, "Erro ao processar atributos livres")
			return
		}
		personagem.AtributosLivres = string(atributosLivresJSON)
	} else {
		personagem.AtributosLivres = "[]"
	}

	// Validar estrutura
	fmt.Printf("DEBUG - CreatePersonagem - Antes da validação: %+v\n", personagem)
	if err := validate.Struct(personagem); err != nil {
		fmt.Printf("DEBUG - CreatePersonagem - Erro na validação: %v\n", err)
		h.Response.BadRequest(c, err.Error())
		return
	}
	fmt.Printf("DEBUG - CreatePersonagem - Validação passou\n")

	// Obtém identificação do usuário e adiciona ao personagem
	sessionID, userIP := middleware.GetUserIdentification(c)
	fmt.Printf("DEBUG - CreatePersonagem - SessionID: %s, UserIP: %s\n", sessionID, userIP)

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

	fmt.Printf("DEBUG - CreatePersonagem - Antes de criar no banco: %+v\n", personagem)

	// Criar no banco
	if err := database.DB.Create(&personagem).Error; err != nil {
		fmt.Printf("DEBUG - CreatePersonagem - Erro ao criar no banco: %v\n", err)
		h.Response.InternalError(c, "Erro ao criar personagem")
		return
	}

	fmt.Printf("DEBUG - CreatePersonagem - Personagem criado com sucesso, ID: %d\n", personagem.ID)

	// Processar perícias selecionadas
	if len(req.PericiasSelecionadas) > 0 {
		// Buscar perícias válidas
		var pericias []models.Pericia
		if err := database.DB.Where("id IN ?", req.PericiasSelecionadas).Find(&pericias).Error; err == nil {
			// Criar registros de PersonagemPericia com fonte 'classe' (padrão para perícias selecionadas na criação)
			for _, pericia := range pericias {
				personagemPericia := models.PersonagemPericia{
					PersonagemID: personagem.ID,
					PericiaID:    pericia.ID,
					Fonte:        "classe",
				}
				if err := database.DB.Create(&personagemPericia).Error; err != nil {
					fmt.Printf("Erro ao associar perícia %d ao personagem: %v\n", pericia.ID, err)
				}
			}
		}
	}

	// Recarregar com relações
	if err := database.DB.Preload("Raca").Preload("Classe").Preload("Origem").Preload("Divindade").First(&personagem, personagem.ID).Error; err != nil {
		h.Response.InternalError(c, "Erro ao carregar personagem criado")
		return
	}

	// Carregar perícias manualmente e calcular stats
	h.loadPersonagemPericias(&personagem)
	h.calculatePersonagemStats(&personagem)

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
	if sessionID != "" && userIP != "" {
		// Se tem ambos, busca por qualquer um dos dois (para lidar com sessões que mudam)
		query = query.Where("id = ? AND (user_session_id = ? OR user_ip = ?)", id, sessionID, userIP)
	} else if sessionID != "" {
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

	var req PersonagemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.Response.BadRequest(c, err.Error())
		return
	}

	// Criar personagem atualizado a partir dos dados da request
	updatedPersonagem := models.Personagem{
		ID:          personagem.ID,
		Nome:        req.Nome,
		Nivel:       req.Nivel,
		For:         req.For,
		Des:         req.Des,
		Con:         req.Con,
		Int:         req.Int,
		Sab:         req.Sab,
		Car:         req.Car,
		RacaID:      req.RacaID,
		ClasseID:    req.ClasseID,
		OrigemID:    req.OrigemID,
		DivindadeID: req.DivindadeID,
	}

	// Garantir que EscolhasRaca seja um JSON válido
	if req.EscolhasRaca == "" {
		updatedPersonagem.EscolhasRaca = "{}"
	} else {
		updatedPersonagem.EscolhasRaca = req.EscolhasRaca
	}

	// LIMPAR relações many2many para evitar problemas no Save
	updatedPersonagem.Pericias = nil

	// Processar atributos livres
	if len(req.AtributosLivres) > 0 {
		atributosLivresJSON, err := json.Marshal(req.AtributosLivres)
		if err != nil {
			h.Response.BadRequest(c, "Erro ao processar atributos livres")
			return
		}
		updatedPersonagem.AtributosLivres = string(atributosLivresJSON)
	} else {
		updatedPersonagem.AtributosLivres = "[]"
	}

	// Restaura identificação original (não deve ser alterada)
	updatedPersonagem.UserSessionID = originalSessionID
	updatedPersonagem.UserIP = originalUserIP
	updatedPersonagem.CreatedByType = originalCreatedByType

	// Usar o personagem atualizado
	personagem = updatedPersonagem

	// Validar estrutura
	if err := validate.Struct(personagem); err != nil {
		h.Response.BadRequest(c, fmt.Sprintf("Erro de validação: %v", err))
		return
	}

	// Salvar alterações
	if err := database.DB.Save(&personagem).Error; err != nil {
		h.Response.InternalError(c, "Erro ao atualizar personagem")
		return
	}

	// Processar perícias selecionadas
	if len(req.PericiasSelecionadas) > 0 {
		// Buscar perícias válidas
		var pericias []models.Pericia
		if err := database.DB.Where("id IN ?", req.PericiasSelecionadas).Find(&pericias).Error; err == nil {
			// Associar perícias ao personagem
			if err := database.DB.Model(&personagem).Association("Pericias").Replace(&pericias); err != nil {
				fmt.Printf("Erro ao associar perícias ao personagem: %v\n", err)
			}
		}
	} else {
		// Se não há perícias selecionadas, limpar associações
		database.DB.Model(&personagem).Association("Pericias").Clear()
	}

	// Carregar perícias manualmente e calcular stats
	h.loadPersonagemPericias(&personagem)
	h.calculatePersonagemStats(&personagem)

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
	if sessionID != "" && userIP != "" {
		// Se tem ambos, busca por qualquer um dos dois (para lidar com sessões que mudam)
		query = query.Where("id = ? AND (user_session_id = ? OR user_ip = ?)", id, sessionID, userIP)
	} else if sessionID != "" {
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

// SavePoderesDivinos salva poderes divinos selecionados para um personagem
func (h *PersonagemHandler) SavePoderesDivinos(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	var request struct {
		PoderesIDs []uint `json:"poderes_ids" validate:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.Response.BadRequest(c, err.Error())
		return
	}

	// Verificar se o personagem existe e pertence ao usuário
	personagem, err := h.findPersonagemByUser(c, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	// Remover poderes divinos existentes
	if err := database.DB.Where("personagem_id = ?", id).Delete(&models.PersonagemPoderDivino{}).Error; err != nil {
		h.Response.InternalError(c, "Erro ao remover poderes existentes")
		return
	}

	// Adicionar novos poderes divinos
	for _, poderID := range request.PoderesIDs {
		poderDivino := models.PersonagemPoderDivino{
			PersonagemID: uint(id),
			PoderID:      poderID,
			Nivel:        personagem.Nivel,
		}

		if err := database.DB.Create(&poderDivino).Error; err != nil {
			h.Response.InternalError(c, "Erro ao salvar poder divino")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Poderes divinos salvos com sucesso",
		"poderes_salvos": len(request.PoderesIDs),
		"personagem_id":  id,
	})
}

// SavePoderesClasse salva poderes de classe selecionados para um personagem
func (h *PersonagemHandler) SavePoderesClasse(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	var request struct {
		PoderesIDs []uint `json:"poderes_ids" validate:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.Response.BadRequest(c, err.Error())
		return
	}

	// Verificar se o personagem existe e pertence ao usuário
	personagem, err := h.findPersonagemByUser(c, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	// Remover poderes de classe existentes
	if err := database.DB.Where("personagem_id = ?", id).Delete(&models.PersonagemPoderClasse{}).Error; err != nil {
		h.Response.InternalError(c, "Erro ao remover poderes existentes")
		return
	}

	// Adicionar novos poderes de classe
	for _, poderID := range request.PoderesIDs {
		poderClasse := models.PersonagemPoderClasse{
			PersonagemID: uint(id),
			PoderID:      poderID,
			Nivel:        personagem.Nivel,
		}

		if err := database.DB.Create(&poderClasse).Error; err != nil {
			h.Response.InternalError(c, "Erro ao salvar poder de classe")
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "Poderes de classe salvos com sucesso",
		"poderes_salvos": len(request.PoderesIDs),
		"personagem_id":  id,
	})
}

// GetAllPersonagensDebug retorna TODOS os personagens do banco (sem filtro de usuário) - APENAS PARA DEBUG
func (h *PersonagemHandler) GetAllPersonagensDebug(c *gin.Context) {
	var personagens []models.Personagem

	// Busca TODOS os personagens sem filtro de usuário
	if err := database.DB.Preload("Raca").Preload("Classe").Preload("Origem").Preload("Divindade").Find(&personagens).Error; err != nil {
		h.Response.InternalError(c, "Erro ao buscar personagens: "+err.Error())
		return
	}

	// Carregar dados completos para cada personagem
	for i := range personagens {
		h.loadPersonagemCompleteData(&personagens[i])
	}

	// Adiciona informações de debug
	result := gin.H{
		"total_personagens": len(personagens),
		"personagens":       personagens,
		"debug_info": gin.H{
			"warning": "Este endpoint mostra TODOS os personagens do banco (sem filtro de usuário)",
			"uso":     "Apenas para debug/administração",
		},
	}

	c.JSON(http.StatusOK, result)
}

// loadPersonagemPericias carrega as perícias de um personagem manualmente
func (h *PersonagemHandler) loadPersonagemPericias(personagem *models.Personagem) {
	var pericias []models.Pericia

	// Query para buscar perícias através da tabela intermediária
	err := h.DB.Table("pericias").
		Select("pericias.*").
		Joins("JOIN personagem_pericias ON pericias.id = personagem_pericias.pericia_id").
		Where("personagem_pericias.personagem_id = ?", personagem.ID).
		Find(&pericias).Error

	if err != nil {
		// Em caso de erro, deixa o slice vazio
		personagem.Pericias = []models.Pericia{}
		return
	}

	personagem.Pericias = pericias
}

// calculatePersonagemStats calcula e preenche os stats (PV, PM, Defesa) de um personagem
func (h *PersonagemHandler) calculatePersonagemStats(personagem *models.Personagem) {
	// Verificar se a Classe está carregada
	if personagem.ClasseID == 0 {
		return // Não pode calcular sem classe
	}

	// Modificadores dos atributos (Tormenta20: modificador = valor do atributo)
	modCon := personagem.Con
	modDes := personagem.Des

	// Buscar dados da classe se não estiver carregada
	var pvPorNivel, pmPorNivel int
	if personagem.Classe.ID == 0 {
		// Buscar os dados da classe
		var classe models.Classe
		if err := h.DB.First(&classe, personagem.ClasseID).Error; err == nil {
			pvPorNivel = classe.PVPorNivel
			pmPorNivel = classe.PMPorNivel
		} else {
			// Valores padrão se não conseguir buscar a classe
			pvPorNivel = 4 // Valor padrão
			pmPorNivel = 2 // Valor padrão
		}
	} else {
		pvPorNivel = personagem.Classe.PVPorNivel
		pmPorNivel = personagem.Classe.PMPorNivel
	}

	// Calcular PV Total = (PV por nível da classe * nível) + (modificador CON * nível)
	pvTotal := (pvPorNivel * personagem.Nivel) + (modCon * personagem.Nivel)

	// Calcular PM Total = PM por nível da classe * nível
	pmTotal := pmPorNivel * personagem.Nivel

	// Calcular Defesa = 10 + modificador DES
	defesa := 10 + modDes

	// Aplicar valores calculados
	personagem.PVTotal = pvTotal
	personagem.PMTotal = pmTotal
	personagem.Defesa = defesa
}

// loadPersonagemCompleteData carrega todas as relações necessárias de um personagem
func (h *PersonagemHandler) loadPersonagemCompleteData(personagem *models.Personagem) {
	// Carregar perícias manualmente (já implementado)
	h.loadPersonagemPericias(personagem)

	// Carregar habilidades e perícias da origem
	if personagem.OrigemID != 0 {
		var origem models.Origem
		if err := h.DB.Preload("Habilidades").Preload("Pericias").First(&origem, personagem.OrigemID).Error; err == nil {
			personagem.Origem = origem
		}
	}

	// Carregar habilidades da divindade
	if personagem.DivindadeID != nil && *personagem.DivindadeID != 0 {
		var divindade models.Divindade
		if err := h.DB.Preload("Habilidades").First(&divindade, *personagem.DivindadeID).Error; err == nil {
			personagem.Divindade = &divindade
		}
	}

	// Carregar habilidades da raça
	if personagem.RacaID != 0 {
		var raca models.Raca
		if err := h.DB.Preload("Habilidades").Preload("Pericias").First(&raca, personagem.RacaID).Error; err == nil {
			personagem.Raca = raca
		}
	}

	// Carregar habilidades da classe
	if personagem.ClasseID != 0 {
		var classe models.Classe
		if err := h.DB.Preload("Habilidades").Preload("PericiasDisponiveis").Preload("PericiasAutomaticas").First(&classe, personagem.ClasseID).Error; err == nil {
			personagem.Classe = classe
		}
	}

	// Calcular stats
	h.calculatePersonagemStats(personagem)
}

// GetPoderesDivinos retorna os poderes divinos selecionados de um personagem
func (h *PersonagemHandler) GetPoderesDivinos(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	// Verificar se o personagem existe e pertence ao usuário
	_, err = h.findPersonagemByUser(c, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	// Buscar poderes divinos selecionados
	var poderesPersonagem []models.PersonagemPoderDivino
	if err := database.DB.Where("personagem_id = ?", id).Find(&poderesPersonagem).Error; err != nil {
		h.Response.InternalError(c, "Erro ao buscar poderes divinos")
		return
	}

	// Extrair apenas os IDs dos poderes
	var poderesIDs []uint
	for _, pp := range poderesPersonagem {
		poderesIDs = append(poderesIDs, pp.PoderID)
	}

	c.JSON(http.StatusOK, gin.H{
		"personagem_id": id,
		"poderes_ids":   poderesIDs,
	})
}

// GetPoderesClasse retorna os poderes de classe selecionados de um personagem
func (h *PersonagemHandler) GetPoderesClasse(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	// Verificar se o personagem existe e pertence ao usuário
	_, err = h.findPersonagemByUser(c, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	// Buscar poderes de classe selecionados
	var poderesPersonagem []models.PersonagemPoderClasse
	if err := database.DB.Where("personagem_id = ?", id).Find(&poderesPersonagem).Error; err != nil {
		h.Response.InternalError(c, "Erro ao buscar poderes de classe")
		return
	}

	// Extrair apenas os IDs dos poderes
	var poderesIDs []uint
	for _, pp := range poderesPersonagem {
		poderesIDs = append(poderesIDs, pp.PoderID)
	}

	c.JSON(http.StatusOK, gin.H{
		"personagem_id": id,
		"poderes_ids":   poderesIDs,
	})
}

// SaveEscolhasRaca salva as escolhas raciais especiais de um personagem
func (h *PersonagemHandler) SaveEscolhasRaca(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	// Parse do body
	var request struct {
		Escolhas map[string]interface{} `json:"escolhas"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		h.Response.BadRequest(c, "Dados inválidos: "+err.Error())
		return
	}

	// Verificar se o personagem existe e pertence ao usuário
	personagem, err := h.findPersonagemByUser(c, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	// Converter as escolhas para JSON
	escolhasJSON, err := json.Marshal(request.Escolhas)
	if err != nil {
		h.Response.InternalError(c, "Erro ao processar escolhas raciais")
		return
	}

	// Atualizar o campo escolhas_raca
	if err := database.DB.Model(&personagem).Update("escolhas_raca", string(escolhasJSON)).Error; err != nil {
		h.Response.InternalError(c, "Erro ao salvar escolhas raciais")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Escolhas raciais salvas com sucesso",
		"personagem_id": id,
	})
}

// GetEscolhasRaca retorna as escolhas raciais especiais de um personagem
func (h *PersonagemHandler) GetEscolhasRaca(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		h.Response.BadRequest(c, "ID inválido")
		return
	}

	// Verificar se o personagem existe e pertence ao usuário
	personagem, err := h.findPersonagemByUser(c, int(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			h.Response.NotFound(c, "Personagem não encontrado")
		} else {
			h.Response.InternalError(c, "Erro ao buscar personagem")
		}
		return
	}

	// Parse das escolhas raciais
	var escolhas map[string]interface{}
	if personagem.EscolhasRaca != "" && personagem.EscolhasRaca != "{}" {
		if err := json.Unmarshal([]byte(personagem.EscolhasRaca), &escolhas); err != nil {
			h.Response.InternalError(c, "Erro ao processar escolhas raciais")
			return
		}
	} else {
		escolhas = make(map[string]interface{})
	}

	c.JSON(http.StatusOK, gin.H{
		"personagem_id": id,
		"escolhas":      escolhas,
	})
}
