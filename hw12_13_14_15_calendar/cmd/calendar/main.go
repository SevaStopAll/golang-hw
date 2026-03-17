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
	internalhttp "github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/server/http"
	"github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage"
	memory "github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage/memory"
	sql "github.com/Ilya19871986/hw-test/hw12_13_14_15_16_calendar/internal/storage/sql"
)

var (
	configFile  string
	showVersion bool
)

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
}

func main() {
	flag.Parse()

	if showVersion {
		printVersion()
		os.Exit(0)
	}

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

	calendarApp := app.New(logg, store)

	httpServer := internalhttp.NewServer(calendarApp, cfg.Server.Host, cfg.Server.Port)

	// Запускаем сервер в горутине для graceful shutdown.
	go func() {
		logg.Info("Starting calendar service...")
		if err := httpServer.Start(); err != nil {
			logg.Error("Failed to start server: " + err.Error())
			os.Exit(1)
		}
	}()

	// Ожидаем сигналы для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logg.Info("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := httpServer.Stop(ctx); err != nil {
		logg.Error("Failed to stop server gracefully: " + err.Error())
	}

	logg.Info("Server stopped")
}
