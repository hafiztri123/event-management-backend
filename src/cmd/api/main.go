package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/hafiztri123/docs"
	"github.com/hafiztri123/src/internal/delivery/handler"
	handlerImplementation "github.com/hafiztri123/src/internal/delivery/handler/implementation"
	"github.com/hafiztri123/src/internal/pkg/cache"
	"github.com/hafiztri123/src/internal/pkg/config"
	"github.com/hafiztri123/src/internal/pkg/database"
	customMiddleware "github.com/hafiztri123/src/internal/pkg/middleware"
	"github.com/hafiztri123/src/internal/repository/implementation"
	"github.com/hafiztri123/src/internal/repository/postgres"
	serviceImplementation "github.com/hafiztri123/src/internal/service/implementation"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/gorm"
)

// @title           Event Management API
// @version         1.0
// @description     API Server for Event Management Application
// @termsOfService  http://swagger.io/terms/

// @contact.name   Hafizh tri Wahyu Muhammad
// @contact.email  hafiz.triwahyu@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.


func main(){
	cfg := loadConfig()
	db := loadDatabase(&cfg.Database)
	startMigration(db)
	redisClient := redisClientInit(cfg.Redis)
	redisCache := redisCacheInit(redisClient, cfg.Redis)
	rateLimitMiddleware := redisRateLimitMiddleware(redisClient, cfg.RateLimit)
	router := chi.NewRouter()
	authMiddleware := customMiddleware.NewAuthMiddleware(cfg.Auth.JWTSecret)
	applyMiddleware(router)
	authHandler := authHandlerInit(db, cfg)
	eventHandler := eventHandlerInit(db, redisCache)
	authRouteInit(authHandler, router)
	healthRouteInit(router)
	eventRouteInit(eventHandler, router, *authMiddleware, *rateLimitMiddleware)
	swaggerRouteInit(router)
	startServer(router)
}

func applyMiddleware(router *chi.Mux) {
	log.Println("[OK] Apply Middleware")
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(60 * time.Second))

	
}

func healthRouteInit(router *chi.Mux) {
	log.Println("[OK] Health route initialization")
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
}

func startServer(router *chi.Mux) {
	log.Println("[OK] Starting server...")
	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig:= make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func(){
		<-sig
		log.Println("[OK] Shuting Down...")

		shutdownCtx,_ := context.WithTimeout(serverCtx, 30*time.Second)

		go func ()  {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded{
				log.Fatal("[FAIL] graceful shutdown timed out... forcing exit")
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed{
		log.Fatal(err)
	}

	<-serverCtx.Done()
}

func loadConfig() *config.Config {
	log.Println("[OK] load config")
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("[FAIL] failed to load config")
	}
	return cfg
}

func loadDatabase(cfg *config.DatabaseConfig) *gorm.DB {
	log.Println("[OK] load database")
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("[FAIL] fail to load database")
	}
	return db
}

func startMigration(db *gorm.DB) {
	log.Println("[OK] start migration")
	postgres.RunMigrations(db)
}

func authHandlerInit(db *gorm.DB, cfg *config.Config) handler.AuthHandler {
	log.Println("[OK] auth handler initialization")
	userRepo :=  repositoryImplementation.NewUserRepository(db)
	authService := serviceImplementation.NewAuthService(userRepo, &cfg.Auth)
	authHandler := handlerImplementation.NewAuthHandler(authService)
	return authHandler
}

func authRouteInit(authHandler handler.AuthHandler, router *chi.Mux) {
	log.Println("[OK] authentication route initialization")
	router.Post("/api/v1/auth/register", authHandler.Register)
	router.Post("/api/v1/auth/login", authHandler.Login)
}


func eventHandlerInit(db *gorm.DB, cfg *cache.RedisCache) handler.EventHandler {
	log.Println("[OK] event handler initialization")
	eventRepo := repositoryImplementation.NewEventRepository(db, cfg)
	eventService := serviceImplementation.NewEventService(eventRepo)
	eventHandler := handlerImplementation.NewEventHandler(eventService)
	return eventHandler
}

func eventRouteInit(eventHandler handler.EventHandler, router *chi.Mux, authMiddleware customMiddleware.AuthMiddleware, rateLimitMiddleware customMiddleware.RateLimiter) {
	log.Println("[OK] event route initialization")
	router.Group(func(r chi.Router) {
		r.Get("/api/v1/events", eventHandler.ListEvents)
		r.Get("/api/v1/events/search", eventHandler.SearchEvents)
		r.Get("/api/v1/events/{id}", eventHandler.GetEvent)
	})

	//Protected
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)
		r.Use(rateLimitMiddleware.RateLimit)
		r.Post("/api/v1/events", eventHandler.CreateEvent)
        r.Put("/api/v1/events/{id}", eventHandler.UpdateEvent)
        r.Delete("/api/v1/events/{id}", eventHandler.DeleteEvent)
	})
}

func swaggerRouteInit(router *chi.Mux) {
	log.Println("[OK] swagger route initialization")
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
}

func redisClientInit(cfg config.RedisConfig) *redis.Client {
	log.Println("[OK] Initializing redis client")
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB: 0,
	})
}

func redisCacheInit(client *redis.Client, cfg config.RedisConfig) *cache.RedisCache {
	log.Println("[OK] Initializing redis cache")
	return cache.NewRedisCache(client, cfg)
}

func redisRateLimitMiddleware(client *redis.Client, cfg config.RateLimitConfig) *customMiddleware.RateLimiter {
	log.Println("[OK] Initializing redis rate limiter middleware")
	return customMiddleware.NewRateLimiter(
		client,
		cfg.RequestLimit,
		time.Duration(cfg.WindowSeconds),
	)
}



