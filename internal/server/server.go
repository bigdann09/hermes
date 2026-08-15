package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bigdann09/notifications/internal/config"
	"github.com/bigdann09/notifications/internal/routes"
	"github.com/bigdann09/notifications/internal/services"
	"github.com/bigdann09/notifications/pkgs/logger"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Server struct {
	logger   *zap.Logger
	config   *config.Config
	server   *http.Server
	engine   *gin.Engine
	services *services.Service
}

func NewServer(cfg *config.Config) *Server {
	logger := logger.NewLogger(&cfg.App)
	logger.Info("setup application services and routes")
	services := services.NewService(logger, cfg)

	router := gin.Default()
	routes := routes.NewRoute(services, router, logger, cfg)
	router = routes.Register()

	logger.Info("configuring server...")
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.App.Port),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	return &Server{
		logger:   logger,
		config:   cfg,
		engine:   router,
		server:   server,
		services: services,
	}
}

func (app *Server) Start() {
	go func() {
		app.logger.Info(
			"starting server",
			zap.String("port", app.config.App.Port),
		)
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			app.logger.Error("server failed to start", zap.Error(err))
		}
	}()

	app.Shutdown()
}

func (app *Server) Shutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.logger.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.server.Shutdown(ctx); err != nil {
		app.logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	db, _ := app.services.Database.DB()
	if err := db.Close(); err != nil {
		app.logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	if err := app.services.Cache.Close(); err != nil {
		app.logger.Fatal("server forced to shutdown", zap.Error(err))
	}

	if err := app.services.KafkaProducer.Close(); err != nil {
		app.logger.Error("failed to close kafka producer", zap.Error(err))
	}

	<-ctx.Done()
	app.logger.Info("server exited")
}
