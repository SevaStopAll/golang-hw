package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/app"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/config"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/logger"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/mq/kafka"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage/notifications"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/storer.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logg, err := logger.NewLogger(cfg.Logger.Level)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Создаем хранилище для уведомлений
	notificationStore, err := notifications.NewPostgresNotificationStorage(cfg.Storage.DSN)
	if err != nil {
		logg.Fatalf("Failed to create notification storage: %v", err)
	}
	defer notificationStore.Close()

	// Создаем и подключаем Kafka consumer с retry
	consumer := kafka.NewConsumer(cfg.Kafka.Brokers, cfg.Kafka.Topic, cfg.Kafka.GroupID)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := consumer.WaitForConnect(ctx, cfg.Kafka.MaxAttempts, cfg.Kafka.RetryBackoff); err != nil {
		logg.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer consumer.Close()

	logg.Info("Successfully connected to Kafka")

	storer := app.NewStorer(notificationStore, consumer, logg)

	// Graceful shutdown
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logg.Info("Shutting down storer...")
		mainCancel()
	}()

	logg.Info("Starting storer...")
	if err := storer.Run(mainCtx); err != nil {
		logg.Errorf("Storer error: %v", err)
		os.Exit(1)
	}
}
