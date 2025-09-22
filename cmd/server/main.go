package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tormenta20-builder/internal/config"
	"tormenta20-builder/internal/database"
	"tormenta20-builder/internal/handlers"
	"tormenta20-builder/internal/middleware"
	"tormenta20-builder/internal/migrations"

	"github.com/gin-gonic/gin"
)

func main() {
	// Carregar configurações
	if err := config.Load(); err != nil {
		log.Fatalf("Erro ao carregar configurações: %v", err)
	}

	// Conectar ao banco de dados
	if err := database.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("Erro ao fechar conexão com banco: %v", err)
		}
	}()

	// Executar migrations
	if err := migrations.RunMigrations(database.DB); err != nil {
		log.Fatalf("Erro ao executar migrations: %v", err)
	}

	// Auto-migrate apenas em desenvolvimento
	if !config.IsProduction() {
		if err := database.AutoMigrate(); err != nil {
			log.Fatalf("Erro ao executar auto-migrate: %v", err)
		}
	}

	// Configurar e iniciar servidor
	server := setupServer()
	if err := startServerWithGracefulShutdown(server); err != nil {
		log.Fatalf("Erro no servidor: %v", err)
	}
}

// setupServer configura o servidor HTTP
func setupServer() *http.Server {
	serverConfig := config.GetServerConfig()

	// Configurar modo do Gin
	gin.SetMode(serverConfig.Mode)

	router := setupRouter()

	return &http.Server{
		Addr:           ":" + serverConfig.Port,
		Handler:        router,
		ReadTimeout:    serverConfig.ReadTimeout,
		WriteTimeout:   serverConfig.WriteTimeout,
		IdleTimeout:    serverConfig.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}
}

// setupRouter configura as rotas e middlewares
func setupRouter() *gin.Engine {
	r := gin.Default()

	// Middleware
	r.Use(middleware.SetupCORS())
	r.Use(middleware.RequestLogger())

	// Health check
	r.GET("/health", healthCheck)

	// API routes
	api := r.Group("/api/v1")
	{
		// Initialize handlers
		racaHandler := handlers.NewRacaHandler()
		classeHandler := handlers.NewClasseHandler()
		origemHandler := handlers.NewOrigemHandler()
		divindadeHandler := handlers.NewDivindadeHandler()
		personagemHandler := handlers.NewPersonagemHandler()
		periciasHandler := handlers.NewPericiasHandler(database.DB)
		habilidadeHandler := handlers.NewHabilidadeHandler()
		poderHandler := handlers.NewPoderHandler()

		// Register routes
		racaHandler.RegisterRoutes(api)
		classeHandler.RegisterRoutes(api)
		origemHandler.RegisterRoutes(api)
		divindadeHandler.RegisterRoutes(api)
		personagemHandler.RegisterRoutes(api)
		habilidadeHandler.RegisterRoutes(api)
		poderHandler.RegisterRoutes(api)

		// Perícias routes
		api.GET("/pericias", periciasHandler.GetPericias)
		api.GET("/pericias/:id", periciasHandler.GetPericia)
		api.GET("/classes/:id/pericias", periciasHandler.GetPericiasClasse)
		api.GET("/racas/:id/pericias", periciasHandler.GetPericiasRaca)
		api.GET("/origens/:id/pericias", periciasHandler.GetPericiasOrigem)
		api.POST("/personagens/:id/pericias", periciasHandler.UpdatePericiasPersonagem)
	}

	return r
}

// healthCheck endpoint de verificação de saúde
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now(),
		"version":   "1.0.0",
	})
}

// startServerWithGracefulShutdown inicia o servidor com graceful shutdown
func startServerWithGracefulShutdown(srv *http.Server) error {
	// Canal para escutar sinais de interrupção
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor em goroutine
	go func() {
		log.Printf("Servidor iniciado na porta %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor: %v", err)
		}
	}()

	// Aguardar sinal de interrupção
	<-quit
	log.Println("Desligando servidor...")

	// Graceful shutdown com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Erro durante graceful shutdown: %v", err)
		return err
	}

	log.Println("Servidor desligado com sucesso")
	return nil
}
