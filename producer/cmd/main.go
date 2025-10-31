package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"producer/internal/config"
	"producer/internal/handler"
	"producer/internal/service"
	"syscall"
	"time"
)

func initLog(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "debug":
		gin.SetMode(gin.DebugMode)
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		gin.SetMode(gin.ReleaseMode)
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		gin.SetMode(gin.DebugMode)
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}

	return logger
}

func main() {
	cfg := config.GetConfig()

	logger := initLog(cfg.Env)
	slog.SetDefault(logger)

	logger.Info("starting")

	msgService := service.New()
	msgHandler := handler.New(msgService)

	router := gin.Default()
	router.POST("/message", msgHandler.Put)

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Address, cfg.Port),
		Handler: router.Handler(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("listen", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutdown Server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Info("Server Shutdown:", err)
	}
	logger.Info("Server exiting")
}
