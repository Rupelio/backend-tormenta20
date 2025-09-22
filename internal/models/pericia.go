package models

import (
	"time"

	"gorm.io/gorm"
)

// Pericia representa uma perícia do Tormenta20
type Pericia struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Nome      string         `json:"nome" gorm:"column:nome;not null;unique"`
	Atributo  string         `json:"atributo" gorm:"column:atributo;not null"` // FOR, DES, CON, INT, SAB, CAR
	Descricao string         `json:"descricao" gorm:"column:descricao;type:text"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// TableName especifica o nome da tabela para o modelo Pericia
func (Pericia) TableName() string {
	return "pericias"
}
