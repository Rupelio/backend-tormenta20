package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var validate = validator.New()

// ResponseHandler padroniza as respostas de erro
type ResponseHandler struct{}

func (rh *ResponseHandler) Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func (rh *ResponseHandler) Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, data)
}

func (rh *ResponseHandler) NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, gin.H{"error": message})
}

func (rh *ResponseHandler) BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, gin.H{"error": message})
}

func (rh *ResponseHandler) InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, gin.H{"error": message})
}

// Função utilitária para parse de ID
func parseID(c *gin.Context) (uint, error) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// GenericService para operações CRUD básicas
type GenericService struct {
	DB       *gorm.DB
	Response *ResponseHandler
}

func NewGenericService(db *gorm.DB) *GenericService {
	return &GenericService{
		DB:       db,
		Response: &ResponseHandler{},
	}
}

// GetAll busca todos os registros de um modelo
func (gs *GenericService) GetAll(c *gin.Context, model interface{}, preload ...string) {
	query := gs.DB
	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.Find(model).Error; err != nil {
		gs.Response.InternalError(c, "Erro ao buscar registros")
		return
	}

	gs.Response.Success(c, model)
}

// GetByID busca um registro por ID
func (gs *GenericService) GetByID(c *gin.Context, model interface{}, notFoundMsg string, preload ...string) {
	id, err := parseID(c)
	if err != nil {
		gs.Response.BadRequest(c, "ID inválido")
		return
	}

	query := gs.DB
	for _, p := range preload {
		query = query.Preload(p)
	}

	if err := query.First(model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			gs.Response.NotFound(c, notFoundMsg)
		} else {
			gs.Response.InternalError(c, "Erro ao buscar registro")
		}
		return
	}

	gs.Response.Success(c, model)
}

// Create cria um novo registro
func (gs *GenericService) Create(c *gin.Context, model interface{}, preload ...string) {
	if err := c.ShouldBindJSON(model); err != nil {
		gs.Response.BadRequest(c, err.Error())
		return
	}

	if err := validate.Struct(model); err != nil {
		gs.Response.BadRequest(c, err.Error())
		return
	}

	if err := gs.DB.Create(model).Error; err != nil {
		gs.Response.InternalError(c, "Erro ao criar registro")
		return
	}

	// Recarregar com relações se especificadas
	if len(preload) > 0 {
		query := gs.DB
		for _, p := range preload {
			query = query.Preload(p)
		}
		query.First(model)
	}

	gs.Response.Created(c, model)
}

// Update atualiza um registro existente
func (gs *GenericService) Update(c *gin.Context, model interface{}, notFoundMsg string, preload ...string) {
	id, err := parseID(c)
	if err != nil {
		gs.Response.BadRequest(c, "ID inválido")
		return
	}

	// Verificar se o registro existe
	if err := gs.DB.First(model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			gs.Response.NotFound(c, notFoundMsg)
		} else {
			gs.Response.InternalError(c, "Erro ao buscar registro")
		}
		return
	}

	// Bind dos novos dados
	if err := c.ShouldBindJSON(model); err != nil {
		gs.Response.BadRequest(c, err.Error())
		return
	}

	// Validar estrutura
	if err := validate.Struct(model); err != nil {
		gs.Response.BadRequest(c, err.Error())
		return
	}

	// Salvar alterações
	if err := gs.DB.Save(model).Error; err != nil {
		gs.Response.InternalError(c, "Erro ao atualizar registro")
		return
	}

	// Recarregar com relações se especificadas
	if len(preload) > 0 {
		query := gs.DB
		for _, p := range preload {
			query = query.Preload(p)
		}
		query.First(model)
	}

	gs.Response.Success(c, model)
}

// Delete remove um registro
func (gs *GenericService) Delete(c *gin.Context, model interface{}, notFoundMsg string) {
	id, err := parseID(c)
	if err != nil {
		gs.Response.BadRequest(c, "ID inválido")
		return
	}

	// Verificar se existe antes de deletar
	if err := gs.DB.First(model, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			gs.Response.NotFound(c, notFoundMsg)
		} else {
			gs.Response.InternalError(c, "Erro ao buscar registro")
		}
		return
	}

	if err := gs.DB.Delete(model, id).Error; err != nil {
		gs.Response.InternalError(c, "Erro ao deletar registro")
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
