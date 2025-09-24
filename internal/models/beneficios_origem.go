package models

// PersonagemBeneficioPericia representa a tabela de junção para perícias de origem
type PersonagemBeneficioPericia struct {
	PersonagemID uint `gorm:"primaryKey"`
	PericiaID    uint `gorm:"primaryKey"`
}

// PersonagemBeneficioPoder representa a tabela de junção para poderes de origem
type PersonagemBeneficioPoder struct {
	PersonagemID uint `gorm:"primaryKey"`
	PoderID      uint `gorm:"primaryKey"`
}
