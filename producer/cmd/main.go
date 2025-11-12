package main

import (
	"context"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"producer/internal/config"
	"producer/internal/handler"
	"producer/internal/queue"
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
		gin.SetMode(gin.DebugMode)

		if err := os.MkdirAll("logs", os.ModeDir|os.ModePerm); err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}

		logFile, err := os.OpenFile("logs/producer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("Failed to open log file: %v", err)
		}

		w := io.MultiWriter(os.Stdout, logFile)
		logger = slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: slog.LevelDebug}))
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

	kafkaProducer, err := queue.NewKafkaProducer(
		cfg.Kafka.MessageTopic,
		&kafka.ConfigMap{
			"bootstrap.servers": cfg.Kafka.BootstrapServers,
			"acks":              cfg.Kafka.Acks},
		cfg.Kafka.FlushTimeout)
	if err != nil {
		logger.Error("failed create Kafka producer", slog.String("error", err.Error()))
		os.Exit(2)
	}

	msgService := service.New(kafkaProducer)
	msgHandler := handler.New(msgService)

	router := gin.Default()
	router.POST("/books", msgHandler.Post)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HttpServer.Address, cfg.HttpServer.Port),
		Handler:      router.Handler(),
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
		WriteTimeout: cfg.HttpServer.Timeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logger.Error("listen", slog.String("error", err.Error()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown Server")

	kafkaProducer.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Info("server Shutdown:", slog.String("error", err.Error()))
	}
	logger.Info("server exiting")
}
