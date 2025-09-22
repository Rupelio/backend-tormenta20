# Tormenta20 Builder - Backend

Backend da aplicação Tormenta20 Builder construído em Go com Gin Framework.

## 🚀 Funcionalidades

- **API RESTful** para gerenciamento de personagens
- **CRUD completo** para raças, classes, origens e personagens
- **Cálculo automático** de estatísticas de personagens
- **Validação robusta** de dados
- **Migrations automáticas** do banco de dados
- **Graceful shutdown** do servidor
- **Middleware personalizado** para CORS, logging e tratamento de erros

## 🛠️ Tecnologias

- **Go 1.21+**
- **Gin Web Framework**
- **GORM** (ORM)
- **PostgreSQL**
- **Docker** & Docker Compose

## 📁 Estrutura do Projeto

```
backend/
├── cmd/server/          # Ponto de entrada da aplicação
├── internal/
│   ├── config/          # Configurações da aplicação
│   ├── database/        # Configuração e conexão do banco
│   ├── handlers/        # Handlers HTTP (controllers)
│   ├── middleware/      # Middlewares customizados
│   ├── migrations/      # Migrations do banco de dados
│   ├── models/          # Modelos de dados (entities)
│   └── services/        # Lógica de negócio
├── Dockerfile           # Configuração Docker
├── Makefile            # Comandos úteis de desenvolvimento
├── go.mod              # Dependências Go
└── go.sum              # Checksums das dependências
```

## 🔧 Configuração

### Variáveis de Ambiente

Crie um arquivo `.env` na pasta `database/` com as seguintes variáveis:

```env
# Banco de Dados
DB_HOST=localhost
DB_USER=postgres
DB_PASS=sua_senha
DB_NAME=tormenta20
DB_PORT=5432
DB_SSLMODE=disable
DB_TIMEZONE=America/Sao_Paulo

# Configurações do Pool de Conexões
DB_MAX_IDLE_CONNS=10
DB_MAX_OPEN_CONNS=100
DB_CONN_MAX_LIFETIME=3600
DB_LOG_LEVEL=warn

# Servidor
SERVER_PORT=8080
SERVER_READ_TIMEOUT=10
SERVER_WRITE_TIMEOUT=10
SERVER_IDLE_TIMEOUT=30
GIN_MODE=debug
```

### Instalação de Dependências

```bash
# Baixar dependências
make deps

# Ou usando go diretamente
go mod download
```

## 🚀 Executando a Aplicação

### Desenvolvimento

```bash
# Executar diretamente
make run

# Ou com hot reload (requer air)
make dev

# Instalar ferramentas de desenvolvimento
make install-tools
```

### Produção

```bash
# Compilar
make build

# Executar binário
./bin/tormenta20-builder
```

### Docker

```bash
# Construir imagem
make docker-build

# Executar container
make docker-run

# Ou usar docker-compose (na raiz do projeto)
docker-compose up
```

## 🧪 Testes

```bash
# Executar todos os testes
make test

# Testes com cobertura
make test-coverage

# Executar benchmarks
make bench
```

## 🔍 Qualidade de Código

```bash
# Formatar código
make fmt

# Análise estática
make vet

# Linter (requer golangci-lint)
make lint

# Executar todas as verificações
make check

# Verificações antes do commit
make pre-commit
```

## 📊 API Endpoints

### Health Check
- `GET /health` - Verificação de saúde da aplicação

### Raças
- `GET /api/v1/racas` - Listar todas as raças
- `GET /api/v1/racas/:id` - Obter raça por ID

### Classes
- `GET /api/v1/classes` - Listar todas as classes
- `GET /api/v1/classes/:id` - Obter classe por ID

### Origens
- `GET /api/v1/origens` - Listar todas as origens
- `GET /api/v1/origens/:id` - Obter origem por ID

### Personagens
- `POST /api/v1/personagens` - Criar personagem
- `GET /api/v1/personagens/:id` - Obter personagem por ID
- `PUT /api/v1/personagens/:id` - Atualizar personagem
- `DELETE /api/v1/personagens/:id` - Deletar personagem
- `POST /api/v1/personagens/calculate` - Calcular estatísticas

## 🗄️ Banco de Dados

### Migrations

As migrations são executadas automaticamente na inicialização da aplicação. Os arquivos estão em `internal/migrations/files/`.

### Modelos

- **Raca**: Informações das raças disponíveis
- **Classe**: Informações das classes disponíveis
- **Origem**: Informações das origens disponíveis
- **HabilidadeRaca**: Habilidades específicas de raças
- **HabilidadeClasse**: Habilidades específicas de classes
- **Personagem**: Dados dos personagens criados

## 🏗️ Arquitetura

O projeto segue uma arquitetura limpa com separação clara de responsabilidades:

### Handlers (Controllers)
Responsáveis por:
- Validação de entrada
- Chamada dos serviços
- Formatação da resposta

### Services
Contêm a lógica de negócio da aplicação:
- Cálculos de estatísticas
- Validações de domínio
- Regras de negócio

### Models
Definem a estrutura dos dados e mapeamento ORM.

### Database
Gerencia conexão, configuração e operações do banco de dados.

## 📋 Padrões Utilizados

- **DRY (Don't Repeat Yourself)**: Eliminação de código duplicado
- **Generic Services**: Operações CRUD reutilizáveis
- **Dependency Injection**: Facilitando testes e manutenção
- **Error Handling**: Tratamento padronizado de erros
- **Configuration Management**: Configuração centralizada e validada

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Comandos Makefile

Use `make help` para ver todos os comandos disponíveis:

```bash
make help
```

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.
