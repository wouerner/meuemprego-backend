package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	_ "github.com/wouerner/runter-backend/docs"
	"github.com/wouerner/runter-backend/internal/config"
	"github.com/wouerner/runter-backend/internal/database"
	"github.com/wouerner/runter-backend/internal/handler"
	"github.com/wouerner/runter-backend/internal/middleware"
	"github.com/wouerner/runter-backend/internal/repository"
	"github.com/wouerner/runter-backend/internal/service"
)

// @title           Runter Backend API
// @version         1.0
// @description     API RESTful em Go com autenticação JWT, PostgreSQL e Chi Router.
// @termsOfService  http://swagger.io/terms/

// @contact.name    Suporte da API
// @contact.email   suporte@runter.com

// @license.name    MIT
// @license.url     https://opensource.org/licenses/MIT

// @host            localhost:8080
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in              header
// @name            Authorization
// @description     Digite 'Bearer ' seguido pelo seu token JWT obtido no login/registro.
func main() {
	// 1. Carregar configurações da aplicação
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Falha ao carregar configurações: %v", err)
	}

	// 2. Inicializar banco de dados com GORM & Postgres
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("Falha ao conectar no banco de dados: %v", err)
	}

	// 3. Inicializar camadas da arquitetura (Repository -> Service -> Handler)
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg)

	candidateRepo := repository.NewCandidateRepository(db)
	candidateService := service.NewCandidateService(candidateRepo)

	hunterRepo := repository.NewHunterRepository(db)
	hunterService := service.NewHunterService(hunterRepo)

	accessRequestRepo := repository.NewAccessRequestRepository(db)
	accessRequestService := service.NewAccessRequestService(accessRequestRepo, hunterRepo, candidateRepo)

	metricRepo := repository.NewMetricRepository(db)
	metricService := service.NewMetricService(metricRepo)

	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(authService)
	candidateHandler := handler.NewCandidateHandler(candidateService)
	hunterHandler := handler.NewHunterHandler(hunterService)
	accessRequestHandler := handler.NewAccessRequestHandler(accessRequestService)
	metricHandler := handler.NewMetricHandler(metricService)
	healthHandler := handler.NewHealthHandler()

	// 4. Configurar roteador Chi e middlewares globais
	r := chi.NewRouter()

	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(60 * time.Second))

	// Configuração de CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Rota para UI do Swagger em /swagger/index.html
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// 5. Roteamento e Versionamento da API (/api/v1)
	r.Route("/api/v1", func(r chi.Router) {
		// Endpoint público de verificação de saúde da API
		r.Get("/health", healthHandler.Check)

		// Rotas públicas de Autenticação
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", authHandler.Register)
			r.Post("/login", authHandler.Login)
		})

		// Rotas protegidas por autenticação JWT
		r.Group(func(r chi.Router) {
			r.Use(middleware.JWTMiddleware(cfg.JWTSecret))

			r.Route("/users", func(r chi.Router) {
				r.Get("/me", userHandler.GetMe)
			})

			r.Get("/candidates", candidateHandler.List)
			r.Get("/candidates/me", candidateHandler.GetMe)
			r.Put("/candidates/me", candidateHandler.SaveMe)
			r.Patch("/candidates/{id}", candidateHandler.SetApproval)

			r.Get("/hunters", hunterHandler.List)
			r.Get("/hunters/me", hunterHandler.GetMe)
			r.Put("/hunters/me", hunterHandler.SaveMe)
			r.Patch("/hunters/{id}/status", hunterHandler.SetStatus)
			r.Post("/hunters/{id}/contacts", hunterHandler.IncrementContacts)

			r.Get("/access-requests/me", accessRequestHandler.ListMe)
			r.Post("/access-requests", accessRequestHandler.Send)
			r.Patch("/access-requests/{id}", accessRequestHandler.Respond)

			r.Get("/metrics", metricHandler.List)
			r.Post("/metrics", metricHandler.Track)
		})
	})

	// 6. Iniciar servidor HTTP
	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("🚀 Servidor Go rodando na porta %s", cfg.Port)
	log.Printf("📚 Documentação Swagger disponível em: http://localhost:%s/swagger/index.html", cfg.Port)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Erro ao iniciar servidor HTTP: %v", err)
	}
}
