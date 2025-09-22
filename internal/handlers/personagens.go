package handlers

import (
	"net/http"

	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/models"

	"github.com/gin-gonic/gin"
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
	}
}

func (h *PersonagemHandler) GetAllPersonagens(c *gin.Context) {
	var personagens []models.Personagem
	h.GetAll(c, &personagens, "Raca", "Classe", "Origem")
}

func (h *PersonagemHandler) GetPersonagem(c *gin.Context) {
	var personagem models.Personagem
	h.GetByID(c, &personagem, "Personagem não encontrado", "Raca", "Classe", "Origem")
}

func (h *PersonagemHandler) CreatePersonagem(c *gin.Context) {
	var personagem models.Personagem
	h.Create(c, &personagem, "Raca", "Classe", "Origem")
}

func (h *PersonagemHandler) UpdatePersonagem(c *gin.Context) {
	var personagem models.Personagem
	h.Update(c, &personagem, "Personagem não encontrado", "Raca", "Classe", "Origem")
}

func (h *PersonagemHandler) DeletePersonagem(c *gin.Context) {
	var personagem models.Personagem
	h.Delete(c, &personagem, "Personagem não encontrado")
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
