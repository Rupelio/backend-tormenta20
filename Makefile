# Makefile para o projeto Tormenta20 Builder Backend

.PHONY: help build run test clean fmt vet deps migrate

# Variáveis
APP_NAME=tormenta20-builder
BINARY_DIR=./bin
MAIN_PATH=./cmd/server

# Cores para output
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

## help: Mostra esta mensagem de ajuda
help:
	@echo "$(GREEN)Makefile para ${APP_NAME}$(NC)"
	@echo ""
	@echo "$(YELLOW)Comandos disponíveis:$(NC)"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' | sed -e 's/^/ /'

## build: Compila a aplicação
build:
	@echo "$(GREEN)Compilando aplicação...$(NC)"
	@mkdir -p $(BINARY_DIR)
	@go build -o $(BINARY_DIR)/$(APP_NAME) $(MAIN_PATH)
	@echo "$(GREEN)✓ Aplicação compilada em $(BINARY_DIR)/$(APP_NAME)$(NC)"

## run: Executa a aplicação
run:
	@echo "$(GREEN)Executando aplicação...$(NC)"
	@go run $(MAIN_PATH)/main.go

## dev: Executa com air (hot reload) - requer air instalado
dev:
	@echo "$(GREEN)Iniciando desenvolvimento com hot reload...$(NC)"
	@air

## test: Executa todos os testes
test:
	@echo "$(GREEN)Executando testes...$(NC)"
	@go test -v ./...

## test-coverage: Executa testes com cobertura
test-coverage:
	@echo "$(GREEN)Executando testes com cobertura...$(NC)"
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Relatório de cobertura gerado em coverage.html$(NC)"

## bench: Executa benchmarks
bench:
	@echo "$(GREEN)Executando benchmarks...$(NC)"
	@go test -bench=. -benchmem ./...

## fmt: Formata código Go
fmt:
	@echo "$(GREEN)Formatando código...$(NC)"
	@go fmt ./...
	@echo "$(GREEN)✓ Código formatado$(NC)"

## vet: Executa go vet (análise estática)
vet:
	@echo "$(GREEN)Executando análise estática...$(NC)"
	@go vet ./...
	@echo "$(GREEN)✓ Análise concluída$(NC)"

## lint: Executa golangci-lint (requer golangci-lint instalado)
lint:
	@echo "$(GREEN)Executando linter...$(NC)"
	@golangci-lint run ./...

## deps: Baixa dependências
deps:
	@echo "$(GREEN)Baixando dependências...$(NC)"
	@go mod download
	@echo "$(GREEN)✓ Dependências baixadas$(NC)"

## deps-update: Atualiza dependências
deps-update:
	@echo "$(GREEN)Atualizando dependências...$(NC)"
	@go get -u ./...
	@go mod tidy
	@echo "$(GREEN)✓ Dependências atualizadas$(NC)"

## tidy: Limpa go.mod e go.sum
tidy:
	@echo "$(GREEN)Limpando módulos...$(NC)"
	@go mod tidy
	@echo "$(GREEN)✓ Módulos limpos$(NC)"

## migrate: Executa migrations do banco
migrate:
	@echo "$(GREEN)Executando migrations...$(NC)"
	@go run $(MAIN_PATH)/main.go -migrate
	@echo "$(GREEN)✓ Migrations executadas$(NC)"

## clean: Remove arquivos temporários e binários
clean:
	@echo "$(GREEN)Limpando arquivos temporários...$(NC)"
	@rm -rf $(BINARY_DIR)
	@rm -f coverage.out coverage.html
	@go clean
	@echo "$(GREEN)✓ Limpeza concluída$(NC)"

## docker-build: Constrói imagem Docker
docker-build:
	@echo "$(GREEN)Construindo imagem Docker...$(NC)"
	@docker build -t $(APP_NAME):latest .
	@echo "$(GREEN)✓ Imagem Docker construída$(NC)"

## docker-run: Executa container Docker
docker-run:
	@echo "$(GREEN)Executando container Docker...$(NC)"
	@docker run -p 8080:8080 --env-file .env $(APP_NAME):latest

## install-tools: Instala ferramentas de desenvolvimento
install-tools:
	@echo "$(GREEN)Instalando ferramentas de desenvolvimento...$(NC)"
	@go install github.com/cosmtrek/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✓ Ferramentas instaladas$(NC)"

## check: Executa todas as verificações (fmt, vet, test)
check: fmt vet test
	@echo "$(GREEN)✓ Todas as verificações passaram$(NC)"

## pre-commit: Executa verificações antes do commit
pre-commit: fmt vet lint test
	@echo "$(GREEN)✓ Pronto para commit$(NC)"
