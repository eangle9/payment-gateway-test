package initiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"pg/initiator/foundation"
	"pg/initiator/platform"
	persistencedb "pg/internal/constant/persistenceDB"
	"pg/internal/handler/middleware"
	"pg/platform/hlog"
	"syscall"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Initiate
//
//	@title						letspay API
//	@version					1.0
//	@description				This is the letspay api.
//
//	@contact.name				letspay Support Email
//	@contact.url				sample@gmail.com
//	@contact.email				info@letspay.com
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@securityDefinitions.basic	BasicAuth
func Initiate() {
	// Initiate Sample Logger for initial setup
	sampleLogger, err := zap.NewProduction()
	if err != nil {
		fmt.Printf(`{"level":"fatal","msg":"failed to initialize sample logger: %v"}\n`, err)
		os.Exit(1)
	}

	// Initiate Config
	sampleLogger.Info("initializing config")
	foundation.InitConfig()
	sampleLogger.Info("config initialized")
	// Initialize Sentry
	foundation.InitSentry(sampleLogger)
	// Initialize main logger
	sampleLogger.Info("Initializing Logger")
	log := hlog.New(foundation.InitLogger(), hlog.Options{}, nil) // Sentry client handling might need adjustment or removal if not directly supported in hlog options
	log.Info(context.Background(), "logger initialized")

	// Initialize state
	log.Info(context.Background(), "initializing state")
	state := foundation.InitState(log)
	log.Info(context.Background(), "state initialized")

	// Initialize Database
	log.Info(context.Background(), "initializing database")
	pgxConn := foundation.InitDB(viper.GetString("DATABASE_URL"), log)
	log.Info(context.Background(), "database initialized")

	// Handle migrations
	if viper.GetBool("MIGRATION_ACTIVE") {
		log.Info(context.Background(), "initializing migration")
		m := foundation.InitiateMigration(viper.GetString("MIGRATION_PATH"), viper.GetString("DATABASE_URL"), log)
		foundation.UpMigration(m, log)
		log.Info(context.Background(), "migration initialized")
	}

	// Initiate Platform
	log.Info(context.Background(), "initializing platform")
	platformInstance := platform.InitPlatform(log, state)
	log.Info(context.Background(), "initialized platform")

	// Initiate Persistence layer
	log.Info(context.Background(), "initializing persistence layer")
	persistence := InitPersistence(persistencedb.New(pgxConn, log, persistencedb.Options{}), log)
	log.Info(context.Background(), "persistence layer initialized")

	// Initiate Module
	log.Info(context.Background(), "initializing module")
	module := InitModule(persistence, log, platformInstance)
	log.Info(context.Background(), "module initialized")

	// Start Worker
	log.Info(context.Background(), "initializing worker")
	go module.PaymentIntent.StartWorker(context.Background())
	log.Info(context.Background(), "worker initialized")

	// Initiate Handler
	log.Info(context.Background(), "initializing handler")
	handler := InitHandler(module, log, viper.GetDuration("SERVER_TIMEOUT"))
	log.Info(context.Background(), "handler initialized")

	// Initiate Router
	log.Info(context.Background(), "initializing server")
	server := echo.New()
	server.HideBanner = true

	// Middleware
	server.Use(echomiddleware.Recover())
	server.Use(middleware.Logger(log.Named("echo")))
	server.Use(sentryecho.New(sentryecho.Options{
		Repanic: true,
	}))
	server.HTTPErrorHandler = middleware.ErrorHandler
	server.Use(foundation.InitCORS())

	log.Info(context.Background(), "server initialized")

	// Setup routes
	log.Info(context.Background(), "initializing router")
	InitRouter(server.Group("/api"), handler, log, platformInstance.Token,
		persistence, platformInstance)
	log.Info(context.Background(), "router initialized")

	// Configure and start server
	addr := viper.GetString("SERVER_HOST") + ":" + viper.GetString("PORT")
	server.Server.ReadHeaderTimeout = viper.GetDuration("SERVER_READ_HEADER_TIMEOUT")

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info(context.Background(), "shutting down server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Fatal(context.Background(), "server forced to shutdown", zap.Error(err))
		}
		log.Info(context.Background(), "server exiting")
	}()

	log.Info(context.Background(), fmt.Sprintf("server is running on %s", addr))
	if err := server.Start(addr); err != nil && err != http.ErrServerClosed {
		log.Fatal(context.Background(), "server failed to start", zap.Error(err))
	}
	log.Info(context.Background(), "server stopped")
}
