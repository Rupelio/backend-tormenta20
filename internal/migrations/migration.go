package migrations

import (
	"embed"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"
)

//go:embed files/*.sql
var migrationFiles embed.FS

type Migration struct {
	ID       uint   `gorm:"primaryKey"`
	Filename string `gorm:"unique"`
	RunAt    string
}

func RunMigrations(db *gorm.DB) error {
	// Criar tabela de migrations se não existir
	if err := db.AutoMigrate(&Migration{}); err != nil {
		return fmt.Errorf("erro ao criar tabela migrations: %v", err)
	}

	// Ler todos os arquivos de migration
	files, err := migrationFiles.ReadDir("files")
	if err != nil {
		return fmt.Errorf("erro ao ler diretório de migrations: %v", err)
	}

	// Ordenar arquivos por nome
	var sqlFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, file.Name())
		}
	}
	sort.Strings(sqlFiles)

	// Executar migrations que ainda não foram rodadas
	for _, filename := range sqlFiles {
		var existingMigration Migration
		result := db.Where("filename = ?", filename).First(&existingMigration)

		if result.Error == gorm.ErrRecordNotFound {
			// Migration não foi executada ainda
			fmt.Printf("Executando migration: %s\n", filename)

			content, err := migrationFiles.ReadFile(fmt.Sprintf("files/%s", filename))
			if err != nil {
				return fmt.Errorf("erro ao ler arquivo %s: %v", filename, err)
			}

			// Dividir por statements (separados por ;)
			statements := strings.Split(string(content), ";")

			for _, stmt := range statements {
				stmt = strings.TrimSpace(stmt)
				if stmt == "" {
					continue
				}

				if err := db.Exec(stmt).Error; err != nil {
					return fmt.Errorf("erro ao executar migration %s: %v\nStatement: %s", filename, err, stmt)
				}
			}

			// Marcar migration como executada
			migration := Migration{
				Filename: filename,
				RunAt:    time.Now().Format(time.RFC3339),
			}
			if err := db.Create(&migration).Error; err != nil {
				return fmt.Errorf("erro ao registrar migration %s: %v", filename, err)
			}

			fmt.Printf("Migration %s executada com sucesso\n", filename)
		}
	}

	return nil
}
