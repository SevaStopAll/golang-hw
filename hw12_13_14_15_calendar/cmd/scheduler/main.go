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
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage"
	memory "github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage/memory"
	sql "github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/scheduler.yaml", "Path to configuration file")
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

	var store storage.Storage
	if cfg.Storage.Type == "sql" {
		store, err = sql.NewStorage(cfg.Storage.DSN)
		if err != nil {
			logg.Fatalf("Failed to create SQL storage: %v", err)
		}
	} else {
		store = memory.NewStorage()
	}
	defer store.Close()

	// Создаем и подключаем Kafka producer с retry
	producer := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := producer.WaitForConnect(ctx, cfg.Kafka.MaxAttempts, cfg.Kafka.RetryBackoff); err != nil {
		logg.Fatalf("Failed to connect to Kafka: %v", err)
	}
	defer producer.Close()

	logg.Info("Successfully connected to Kafka")

	calendarApp := app.New(logg, store)
	scheduler := app.NewScheduler(calendarApp, producer, logg, cfg.Scheduler)

	// Graceful shutdown
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logg.Info("Shutting down scheduler...")
		mainCancel()
	}()

	logg.Info("Starting scheduler...")
	if err := scheduler.Run(mainCtx); err != nil {
		logg.Errorf("Scheduler error: %v", err)
		os.Exit(1)
	}
}
