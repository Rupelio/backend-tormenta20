package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm/logger"
)

// DatabaseConfig estrutura as configurações do banco de dados
type DatabaseConfig struct {
	Host            string
	User            string
	Password        string
	Name            string
	Port            string
	SSLMode         string
	TimeZone        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	LogLevel        logger.LogLevel
}

// ServerConfig estrutura as configurações do servidor
type ServerConfig struct {
	Port         string
	Mode         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// Config contém todas as configurações da aplicação
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

var appConfig *Config

// Load carrega as configurações da aplicação
func Load() error {
	// Tentar carregar .env file (opcional em produção)
	if err := godotenv.Load("../../database/.env"); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis de ambiente do sistema")
	}

	var err error
	appConfig, err = loadConfig()
	if err != nil {
		return fmt.Errorf("erro ao carregar configurações: %w", err)
	}

	return nil
}

// loadConfig carrega e valida todas as configurações
func loadConfig() (*Config, error) {
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return nil, err
	}

	serverConfig := loadServerConfig()

	return &Config{
		Database: dbConfig,
		Server:   serverConfig,
	}, nil
}

// loadDatabaseConfig carrega configurações do banco de dados
func loadDatabaseConfig() (DatabaseConfig, error) {
	config := DatabaseConfig{
		Host:     getEnvOrDefault("DB_HOST", "localhost"),
		User:     getEnvOrDefault("DB_USER", "postgres"),
		Password: getEnvOrDefault("DB_PASS", ""),
		Name:     getEnvOrDefault("DB_NAME", "tormenta20"),
		Port:     getEnvOrDefault("DB_PORT", "5432"),
		SSLMode:  getEnvOrDefault("DB_SSLMODE", "disable"),
		TimeZone: getEnvOrDefault("DB_TIMEZONE", "America/Sao_Paulo"),
	}

	// Validar campos obrigatórios
	if config.Password == "" {
		return config, fmt.Errorf("DB_PASS é obrigatório")
	}

	// Configurações de connection pool
	maxIdle, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_IDLE_CONNS", "10"))
	maxOpen, _ := strconv.Atoi(getEnvOrDefault("DB_MAX_OPEN_CONNS", "100"))
	maxLifetime, _ := strconv.Atoi(getEnvOrDefault("DB_CONN_MAX_LIFETIME", "3600"))

	config.MaxIdleConns = maxIdle
	config.MaxOpenConns = maxOpen
	config.ConnMaxLifetime = time.Duration(maxLifetime) * time.Second

	// Log level
	logLevel := getEnvOrDefault("DB_LOG_LEVEL", "warn")
	switch logLevel {
	case "silent":
		config.LogLevel = logger.Silent
	case "error":
		config.LogLevel = logger.Error
	case "warn":
		config.LogLevel = logger.Warn
	case "info":
		config.LogLevel = logger.Info
	default:
		config.LogLevel = logger.Warn
	}

	return config, nil
}

// loadServerConfig carrega configurações do servidor
func loadServerConfig() ServerConfig {
	readTimeout, _ := strconv.Atoi(getEnvOrDefault("SERVER_READ_TIMEOUT", "10"))
	writeTimeout, _ := strconv.Atoi(getEnvOrDefault("SERVER_WRITE_TIMEOUT", "10"))
	idleTimeout, _ := strconv.Atoi(getEnvOrDefault("SERVER_IDLE_TIMEOUT", "30"))

	return ServerConfig{
		Port:         getEnvOrDefault("SERVER_PORT", "8080"),
		Mode:         getEnvOrDefault("GIN_MODE", "debug"),
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}
}

// getEnvOrDefault retorna o valor da variável de ambiente ou o padrão
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetDatabaseConfig retorna as configurações do banco de dados
func GetDatabaseConfig() DatabaseConfig {
	if appConfig == nil {
		log.Fatal("Configurações não foram carregadas. Chame config.Load() primeiro.")
	}
	return appConfig.Database
}

// GetServerConfig retorna as configurações do servidor
func GetServerConfig() ServerConfig {
	if appConfig == nil {
		log.Fatal("Configurações não foram carregadas. Chame config.Load() primeiro.")
	}
	return appConfig.Server
}

// Get retorna o valor de uma variável de ambiente
func Get(key string) string {
	return os.Getenv(key)
}

// IsProduction verifica se está em modo de produção
func IsProduction() bool {
	return GetServerConfig().Mode == "release"
}
