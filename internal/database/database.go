package database

import (
	"fmt"
	"log"

	"tormenta20-builder/internal/config"
	"tormenta20-builder/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Connect se conecta ao banco de dados usando configurações
func Connect() error {
	cfg := config.GetDatabaseConfig()

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.Name, cfg.Port, cfg.SSLMode, cfg.TimeZone)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	})

	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	// Configurar connection pool
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("erro ao obter instância SQL DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	log.Println("Conexão com o banco de dados estabelecida com sucesso!")
	return nil
}

// AutoMigrate executa auto-migration dos modelos
func AutoMigrate() error {
	return DB.AutoMigrate(
		&models.Raca{},
		&models.Classe{},
		&models.Origem{},
		&models.Divindade{},
		&models.HabilidadeRaca{},
		&models.HabilidadeClasse{},
		&models.Personagem{},
		&models.PersonagemBeneficioPericia{},
		&models.PersonagemBeneficioPoder{},
	)
}

// Close fecha a conexão com o banco
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
