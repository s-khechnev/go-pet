package main

import (
	"consumer/internal/config"
	bookgrpc "consumer/internal/grpc"
	"consumer/internal/queue"
	"consumer/internal/service/analytics"
	"consumer/internal/service/processor"
	"consumer/internal/storage/postgresql"
	"context"
	"fmt"
	confluentkafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func initLog(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case "debug":
		gin.SetMode(gin.DebugMode)
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "prod":
		gin.SetMode(gin.DebugMode)

		if err := os.MkdirAll("logs", os.ModePerm); err != nil {
			log.Fatalf("Failed to create logs directory: %v", err)
		}

		logFile, err := os.OpenFile("logs/consumer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	bookRepo, err := postgresql.NewStorage(ctx, cfg.DB.ConnString)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	bookProcessorService := processor.NewBookProcessorService(bookRepo)

	kafkaConsumer, err := queue.NewConsumer(
		ctx,
		bookProcessorService,
		cfg.Kafka.MessageTopic,
		cfg.Kafka.GroupId,
		&confluentkafka.ConfigMap{
			"bootstrap.servers":  cfg.Kafka.BootstrapServers,
			"session.timeout.ms": int(cfg.Kafka.SessionTimeout.Milliseconds()),
			"auto.offset.reset":  cfg.Kafka.AutoOffsetReset,
		},
		cfg.Kafka.PollTimeout,
	)
	if err != nil {
		log.Fatalf("failed create kafka consumer: %v", err)
	}

	go kafkaConsumer.Run()

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.GrpcServer.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	analyticsService := analytics.NewBookAnalyticsService(bookRepo)
	serverApi := bookgrpc.NewServerApi(analyticsService)
	bookgrpc.Register(grpcServer, serverApi)

	go func() {
		if err := grpcServer.Serve(l); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("shutdown Server")

	grpcServer.GracefulStop()

	logger.Info("server exiting")
}
