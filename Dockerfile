# backend/Dockerfile
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Instalar dependências
RUN apk add --no-cache git

# Copiar go.mod e go.sum
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fonte
COPY . .

RUN go mod tidy

# Build da aplicação
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Estágio final
FROM alpine:latest

# 1. Instalar pacotes necessários
RUN apk --no-cache add ca-certificates

# 2. Criar um diretório para a aplicação
WORKDIR /app

# 3. Copiar o binário para o diretório de trabalho
COPY --from=builder /app/server .

# 4. Copiar as migrations para um subdiretório
COPY --from=builder /app/internal/migrations/files ./internal/migrations/files

# 5. Expor a porta
EXPOSE 8080

# 6. Rodar a aplicação
CMD ["./server"]

