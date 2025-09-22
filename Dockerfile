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

RUN apk --no-cache add ca-certificates tzdata wget
WORKDIR /root/

# Copiar o binário
COPY --from=builder /app/server .

# Copiar migrations
COPY --from=builder /app/internal/migrations/files ./internal/migrations/files

EXPOSE 8080

CMD ["./server"]
