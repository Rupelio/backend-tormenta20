// internal/models/personagem.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Personagem struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Nome  string `json:"nome" validate:"required,min=2"`
	Nivel int    `json:"nivel" validate:"min=1,max=20"`

	// Atributos base (0 a 4)
	For int `json:"for" gorm:"column:for" validate:"min=0,max=4"`
	Des int `json:"des" gorm:"column:des" validate:"min=0,max=4"`
	Con int `json:"con" gorm:"column:con" validate:"min=0,max=4"`
	Int int `json:"int" gorm:"column:int" validate:"min=0,max=4"`
	Sab int `json:"sab" gorm:"column:sab" validate:"min=0,max=4"`
	Car int `json:"car" gorm:"column:car" validate:"min=0,max=4"`

	// Relações
	RacaID      uint       `json:"raca_id"`
	Raca        Raca       `json:"raca" gorm:"foreignKey:RacaID"`
	ClasseID    uint       `json:"classe_id"`
	Classe      Classe     `json:"classe" gorm:"foreignKey:ClasseID"`
	OrigemID    uint       `json:"origem_id"`
	Origem      Origem     `json:"origem" gorm:"foreignKey:OrigemID"`
	DivindadeID *uint      `json:"divindade_id" gorm:"default:null"`
	Divindade   *Divindade `json:"divindade" gorm:"foreignKey:DivindadeID"`

	// Perícias do personagem
	Pericias []Pericia `json:"pericias" gorm:"many2many:personagem_pericias;"`

	// Escolhas específicas de raça (JSON)
	EscolhasRaca string `json:"escolhas_raca" gorm:"column:escolhas_raca;type:jsonb;default:'{}'"`

	// Identificação do usuário/sessão
	UserSessionID *string `json:"user_session_id" gorm:"column:user_session_id;type:varchar(36)"`
	UserIP        *string `json:"user_ip" gorm:"column:user_ip;type:inet"`
	CreatedByType string  `json:"created_by_type" gorm:"column:created_by_type;default:'session'"`

	// Stats calculados (não salvos no DB)
	PVTotal int `json:"pv_total" gorm:"-"`
	PMTotal int `json:"pm_total" gorm:"-"`
	Defesa  int `json:"defesa" gorm:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Personagem) TableName() string {
	return "personagens"
}

type Raca struct {
	gorm.Model
	Nome string `json:"nome"`

	// Modificadores de atributos (flexível)
	AtributoBonus1 string `json:"atributo_bonus_1" gorm:"column:atributo_bonus_1"`
	ValorBonus1    int    `json:"valor_bonus_1" gorm:"column:valor_bonus_1"`
	AtributoBonus2 string `json:"atributo_bonus_2" gorm:"column:atributo_bonus_2"`
	ValorBonus2    int    `json:"valor_bonus_2" gorm:"column:valor_bonus_2"`
	AtributoBonus3 string `json:"atributo_bonus_3" gorm:"column:atributo_bonus_3"`
	ValorBonus3    int    `json:"valor_bonus_3" gorm:"column:valor_bonus_3"`

	// Suporte a penalidades
	AtributoPenalidade string `json:"atributo_penalidade" gorm:"column:atributo_penalidade"`
	ValorPenalidade    int    `json:"valor_penalidade" gorm:"column:valor_penalidade"`

	// Características físicas
	Tamanho      string `json:"tamanho"`
	Deslocamento int    `json:"deslocamento"`

	// Informações descritivas
	Descricao string `json:"descricao"`

	// Relações
	Habilidades []HabilidadeRaca `json:"habilidades" gorm:"foreignKey:RacaID"`
	Pericias    []Pericia        `json:"pericias" gorm:"many2many:raca_pericias;"`
}

type Classe struct {
	gorm.Model
	Nome                string             `json:"nome"`
	PVPorNivel          int                `json:"pvpornivel" gorm:"column:pv_por_nivel"`
	PMPorNivel          int                `json:"pmpornivel" gorm:"column:pm_por_nivel"`
	AtributoPrincipal   string             `json:"atributoprincipal" gorm:"column:atributo_principal"`
	PericiasQuantidade  int                `json:"pericias_quantidade" gorm:"column:pericias_quantidade;default:2"`
	Habilidades         []HabilidadeClasse `json:"habilidades" gorm:"foreignKey:ClasseID"`
	PericiasDisponiveis []Pericia          `json:"pericias_disponiveis" gorm:"many2many:classe_pericias_disponiveis;"`
	PericiasAutomaticas []Pericia          `json:"pericias_automaticas" gorm:"many2many:classe_pericias_automaticas;"`
}

type Origem struct {
	gorm.Model
	Nome        string             `json:"nome"`
	Descricao   string             `json:"descricao"`
	Pericias    []Pericia          `json:"pericias" gorm:"many2many:origem_pericias;"`
	Habilidades []HabilidadeOrigem `json:"habilidades" gorm:"foreignKey:OrigemID"`
}

func (Origem) TableName() string {
	return "origens"
}

type Divindade struct {
	gorm.Model
	Nome        string                `json:"nome"`
	Descricao   string                `json:"descricao"`
	Dominio     string                `json:"dominio"`
	Alinhamento string                `json:"alinhamento"`
	Habilidades []HabilidadeDivindade `json:"habilidades" gorm:"foreignKey:DivindadeID"`
}

type HabilidadeRaca struct {
	gorm.Model
	RacaID      uint   `json:"raca_id"`
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	Opcional    bool   `json:"opcional"`     // Se o jogador pode escolher esta habilidade
	NivelMinimo int    `json:"nivel_minimo"` // Nível mínimo para obter esta habilidade
}

type HabilidadeClasse struct {
	gorm.Model
	ClasseID  uint   `json:"classe_id"`
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
	Nivel     int    `json:"nivel"`    // Em que nível a classe ganha esta habilidade
	Opcional  bool   `json:"opcional"` // Se o jogador pode escolher esta habilidade
}

type HabilidadeOrigem struct {
	gorm.Model
	OrigemID  uint   `json:"origem_id"`
	Nome      string `json:"nome"`
	Descricao string `json:"descricao"`
	Opcional  bool   `json:"opcional"`
}

func (HabilidadeOrigem) TableName() string {
	return "habilidade_origens"
}

type HabilidadeDivindade struct {
	gorm.Model
	DivindadeID uint   `json:"divindade_id"`
	Nome        string `json:"nome"`
	Descricao   string `json:"descricao"`
	Nivel       int    `json:"nivel"`    // Nível de devoto necessário
	Opcional    bool   `json:"opcional"` // Se é uma concessão opcional
}

func (HabilidadeDivindade) TableName() string {
	return "habilidade_divindades"
}

type Poder struct {
	gorm.Model
	Nome       string `json:"nome"`
	Descricao  string `json:"descricao"`
	Tipo       string `json:"tipo"` // Combate, Destino, Magia, Origem
	Requisitos string `json:"requisitos"`
}

func (Poder) TableName() string {
	return "poderes"
}

type RacaHabilidadeEspecial struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	RacaID    uint      `json:"raca_id"`
	Nome      string    `json:"nome"`
	Descricao string    `json:"descricao"`
	Tipo      string    `json:"tipo"`   // versatilidade, deformidade, etc
	Opcoes    string    `json:"opcoes"` // JSON com opções específicas
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (RacaHabilidadeEspecial) TableName() string {
	return "raca_habilidades_especiais"
}

// PersonagemPericia representa a tabela intermediária com fonte
type PersonagemPericia struct {
	PersonagemID uint   `json:"personagem_id" gorm:"primaryKey"`
	PericiaID    uint   `json:"pericia_id" gorm:"primaryKey"`
	Fonte        string `json:"fonte" gorm:"primaryKey;type:varchar(50)"`
}

func (PersonagemPericia) TableName() string {
	return "personagem_pericias"
}
