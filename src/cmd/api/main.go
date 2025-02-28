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
	"github.com/hafiztri123/src/internal/pkg/cache"
	"github.com/hafiztri123/src/internal/pkg/config"
	"github.com/hafiztri123/src/internal/pkg/database"
	"github.com/hafiztri123/src/internal/pkg/health"
	"github.com/hafiztri123/src/internal/pkg/logger"
	customMiddleware "github.com/hafiztri123/src/internal/pkg/middleware"
	"github.com/hafiztri123/src/internal/pkg/storage"
	"github.com/hafiztri123/src/internal/repository"
	"github.com/hafiztri123/src/internal/repository/postgres"
	"github.com/hafiztri123/src/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
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
	appLogger, err := initLogger()
	if err != nil {
		log.Fatalf("[FAIL] Failed to initialize logger: %v", err)
	}
	defer appLogger.Close()


	ctx := context.Background()

	appLogger.Info(ctx, "Application started successfully", map[string]interface{}{
	    "test": true,
	    "timestamp": time.Now().String(),
	})

	cfg := loadConfig(appLogger, ctx)

	db := loadDatabase(appLogger, ctx, &cfg.Database)
	startMigration(appLogger, ctx, db)

	redisClient := redisClientInit(appLogger, ctx, cfg.Redis)
	redisCache := redisCacheInit(appLogger, ctx, redisClient, cfg.Redis)

	repository := newMainRepository(db, redisCache)
	service := newMainService(repository, cfg)
	handler := newMainHandler(service)
	middleware := newMainMiddleware(cfg, redisClient)



	router := chi.NewRouter()
	applyMiddleware(appLogger, router, middleware)

	mainRoute := newMainRoute(appLogger, ctx, handler, router, middleware)
	mainRoute.auth()
	mainRoute.event()
	mainRoute.user()
	mainRoute.category()

	sideRoute := newSideRoute(appLogger, ctx, router, db, redisClient)
	sideRoute.health()
	sideRoute.swagger()

	startServer(appLogger, ctx, router)
}

type mainRoute struct {
	auth func()
	event func()
	user func()
	category func()
}

type sideRoute struct {
	health func()
	swagger func()
}

func newMainRoute(log *logger.Logger, ctx context.Context, handler *mainHandler, router *chi.Mux, middleware *mainMiddleware) *mainRoute {
	return &mainRoute{
		auth: authRouteInit(log, ctx, handler.Auth, router ),
		event: eventRouteInit(log, ctx, handler.Event, router, middleware.JWT, middleware.RateLimiter),
		user: userRouteInit(log, ctx, handler.User, router, middleware.JWT),
		category: categoryRouteInit(log, ctx, handler.Category, router, middleware),
	}
}

func newSideRoute(log *logger.Logger, ctx context.Context, router *chi.Mux, db *gorm.DB, redis *redis.Client) *sideRoute {
	return &sideRoute{
		health: healthRouteInit(log, ctx, router, db, redis),
		swagger: swaggerRouteInit(log, ctx, router),
	}
}


func initLogger() (*logger.Logger, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	// Ensure log directory exists
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, fmt.Errorf("failed to create logs directory: %w", err)
	}

	config := logger.Config{
		AppName:       "event-management-api",
		Environment:   env,
		MinLogLevel:   logger.InfoLevel,
		EnableConsole: true,
		EnableFile:    true,
		LogFilePath:   "/home/hafizh/projects/ongoing/event-management-backend/logs/application.log",
	}

	return logger.New(config)
}

func applyMiddleware(log *logger.Logger, router *chi.Mux, customMiddleware *mainMiddleware) {
	router.Use(corsMiddleware())
	log.Info(context.Background(), "Applying middleware", nil)
	router.Use(customMiddleware.Logger.LoggerMiddleware(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Timeout(60 * time.Second))
}

func healthRouteInit(log *logger.Logger, ctx context.Context, router *chi.Mux, db *gorm.DB, redis *redis.Client) func() {
	return func ()  {
		log.Info(ctx, "Initializing health routes", nil)

		healthChecker := health.NewHealthChecker("1.0.0")
		healthChecker.AddChecker(health.NewDatabaseChecker(db))
		healthChecker.AddChecker(health.NewRedisChecker(redis))
		healthChecker.AddChecker(health.NewMemoryChecker())
		healthChecker.AddChecker(health.NewDiskChecker("."))
	
		router.Get("/api/v1/health", healthChecker.Handler())
	
		router.Get("/api/v1/health/liveness", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		})
	
		router.Get("/api/v1/health/readiness", healthChecker.Handler())		
	}
}

func startServer(log *logger.Logger, ctx context.Context, router *chi.Mux) {
	log.Info(ctx, "Starting server", map[string]interface{}{"port": 8080})

	server := &http.Server{
		Addr: ":8080",
		Handler: router,
	}

	serverCtx, serverStopCtx := context.WithCancel(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig
		shutdownCtx, cancel := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancel()

		log.Info(shutdownCtx, "Shutting down server", nil)

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal(context.Background(), "Graceful shutdown timed out, forcing exit", nil, nil)
			}
		}()

		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal(shutdownCtx, "Server shutdown failed", err, nil)
		}

		serverStopCtx()
	}()

	log.Info(ctx, "Server is ready to handle requests", nil)
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(ctx, "Server failed", err, nil)
	}

	<-serverCtx.Done()
	log.Info(ctx, "Server stopped", nil)
}

func loadConfig(log *logger.Logger, ctx context.Context) *config.Config {
	log.Info(ctx, "Loading configuration", nil)
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal(ctx, "Failed to load configuration", err, nil)
	}
	return cfg
}

func loadDatabase(log *logger.Logger, ctx context.Context, cfg *config.DatabaseConfig) *gorm.DB {
	log.Info(ctx, "Connecting to database", map[string]interface{}{
		"host": cfg.Host,
		"port": cfg.Port,
		"name": cfg.DBName,
	})

	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal(ctx, "Failed to connect to database", err, nil)
	}

	log.Info(ctx, "Connected to database successfully", nil)
	return db
}

func startMigration(log *logger.Logger, ctx context.Context, db *gorm.DB) {
	log.Info(ctx, "Starting database migrations", nil)
	postgres.RunMigrations(db)
	log.Info(ctx, "Database migrations completed", nil)
}


func authRouteInit(log *logger.Logger, ctx context.Context, authHandler handler.AuthHandler, router *chi.Mux) func() {
	return func ()  {
		log.Info(ctx, "Initializing authentication routes", nil)
		router.Post("/api/v1/auth/register", authHandler.Register)
		router.Post("/api/v1/auth/login", authHandler.Login)	
		
	}
}


func eventRouteInit(log *logger.Logger, ctx context.Context, eventHandler handler.EventHandler, router *chi.Mux, authMiddleware *customMiddleware.AuthMiddleware, rateLimitMiddleware *customMiddleware.RateLimiter) func() {
	return func ()  {
		log.Info(ctx, "Initializing event routes", nil)
		router.Group(func(r chi.Router) {
			r.Get("/api/v1/events", eventHandler.ListEvents)
			r.Get("/api/v1/events/search", eventHandler.SearchEvents)
			r.Get("/api/v1/events/{id}", eventHandler.GetEvent)
		})

		router.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Use(rateLimitMiddleware.RateLimit)
			r.Post("/api/v1/events", eventHandler.CreateEvent)
			r.Put("/api/v1/events/{id}", eventHandler.UpdateEvent)
			r.Delete("/api/v1/events/{id}", eventHandler.DeleteEvent)
			r.Post("/api/v1/events/{id}/upload", eventHandler.UploadFile)
		})
	
	}

}

func swaggerRouteInit(log *logger.Logger, ctx context.Context, router *chi.Mux) func() {
	return func ()  {
		log.Info(ctx, "Initializing swagger routes", nil)
		router.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		))
	}
}

func redisClientInit(log *logger.Logger, ctx context.Context, cfg config.RedisConfig) *redis.Client {
	log.Info(ctx, "Initializing Redis client", map[string]interface{}{
		"host": cfg.Host,
		"port": cfg.Port,
	})

	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB: 0,
	})
}

func redisCacheInit(log *logger.Logger, ctx context.Context, client *redis.Client, cfg config.RedisConfig) *cache.RedisCache {
	log.Info(ctx, "Initializing Redis cache", map[string]interface{}{
		"expiration_minutes": cfg.DurationMinute,
	})
	return cache.NewRedisCache(client, cfg)
}




func userRouteInit(log *logger.Logger, ctx context.Context, userHandler handler.UserHandler, router *chi.Mux, authMiddleware *customMiddleware.AuthMiddleware) func() {
	return func ()  {
		log.Info(ctx, "Initializing user routes", nil)

		router.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)
			r.Put("/api/v1/users/profile", userHandler.UpdateProfile)
			r.Get("/api/v1/users/profile", userHandler.GetProfile)
			r.Put("/api/v1/users/password", userHandler.ChangePassword)
			r.Put("/api/v1/users/profile-image", userHandler.UploadProfileImage)
		})	
	}
}

func categoryRouteInit(log *logger.Logger, ctx context.Context, categoryHandler handler.CategoryHandler, router *chi.Mux, middleware *mainMiddleware) func() {
	return func ()  {
		log.Info(ctx, "Initializing category routes", nil)

		router.Group(func(r chi.Router) {
			r.Get("/api/v1/categories", categoryHandler.ListCategories)
			r.Get("/api/v1/categories/{id}", categoryHandler.GetCategory)
		})
	
		router.Group(func(r chi.Router) {
			r.Use(middleware.JWT.Authenticate)
			r.Use(middleware.RateLimiter.RateLimit)
			r.Post("/api/v1/categories", categoryHandler.CreateCategory)
			r.Put("/api/v1/categories/{id}", categoryHandler.UpdateCategory)
			r.Delete("/api/v1/categories/{id}", categoryHandler.DeleteCategory)
		})	
	}
}

type mainRepository struct {
	User 		repository.UserRepository
	Event 		repository.EventRepository
	Category 	repository.CategoryRepository
}

func newMainRepository(db *gorm.DB, cache *cache.RedisCache) *mainRepository {
	return &mainRepository{
		User: 		repository.NewUserRepository(db),
		Event: 		repository.NewEventRepository(db, cache),
		Category: 	repository.NewCategoryRepository(db, cache),
	}
}


type mainService struct {
	Auth 		service.AuthService
	User 		service.UserService
	Category 	service.CategoryService
	Event 		service.EventService
}

func newMainService (repository *mainRepository, cfg *config.Config) *mainService {

	cloudinary, err := storage.NewCloudinaryService(cfg)
	if err != nil {
		log.Fatal("Cloudinary failed")
	}
	
	return &mainService{
		Auth: 		service.NewAuthService(repository.User, &cfg.Auth),
		User: 		service.NewUserService(repository.User, cloudinary),
		Category: 	service.NewCategoryService(repository.Category),
		Event: 		service.NewEventService(repository.Event, repository.Category, cloudinary),
	}
}

type mainHandler struct {
	Auth 		handler.AuthHandler
	User	 	handler.UserHandler
	Category 	handler.CategoryHandler
	Event 		handler.EventHandler
}

func newMainHandler (service *mainService) *mainHandler {
	return &mainHandler{
		Auth: 		handler.NewAuthHandler(service.Auth),
		User: 		handler.NewUserHandler(service.User),
		Category: 	handler.NewCategoryHandler(service.Category),
		Event: 		handler.NewEventHandler(service.Event),
	}
}

type mainMiddleware struct {
	JWT 		*customMiddleware.AuthMiddleware
	RateLimiter *customMiddleware.RateLimiter
	Logger 		*customMiddleware.Logger

}

func newMainMiddleware(cfg *config.Config, redis *redis.Client) *mainMiddleware {
	JWT := customMiddleware.NewAuthMiddleware(
		cfg.Auth.JWTSecret,
	)

	RateLimiter := customMiddleware.NewRateLimiter(
		redis, cfg.RateLimit.RequestLimit, time.Duration(cfg.RateLimit.WindowSeconds),
	)
	
	return &mainMiddleware{
		JWT: JWT,
		RateLimiter: RateLimiter,
		Logger: customMiddleware.NewLogger(),
	}
}

func corsMiddleware() func(http.Handler) http.Handler {
    return cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:5173"}, // Replace with your frontend URL(s)
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
        ExposedHeaders:   []string{"Link"},
        AllowCredentials: true,
        MaxAge:           300, 
    }).Handler
}
