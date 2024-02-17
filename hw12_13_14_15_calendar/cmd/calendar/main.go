package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/efedyakov/go-hw/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const (
	storageTypeMemory = "memory"
	storageTypeSQL    = "sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "envconfig", "config.toml", "Path to configuration file")
}

func NewStorage(ctx context.Context, config Config) (app.Storage, error) {
	var storage app.Storage

	switch config.App.StorageType {
	case storageTypeMemory:
		storage = memorystorage.New()
	case storageTypeSQL:
		ss, err := sqlstorage.New(config.PostgreSQL.BuildDSN())
		if err != nil {
			return nil, err
		}
		if err := ss.Connect(ctx); err != nil {
			return nil, err
		}
		storage = ss
	default:
		return nil, errors.New("storage type not supported")
	}

	return storage, nil
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	if err := godotenv.Load(configFile); err != nil {
		panic(err)
	}

	config := NewConfig()
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	logg, _ := logger.New(config.Logger.Level)
	ctx := context.Background()

	storage, err := NewStorage(ctx, config)
	if err != nil {
		logg.Fatal("failed to create a storage: " + err.Error())
	}
	if s, ok := storage.(sqlstorage.Storage); ok {
		defer s.Close(ctx)
	}

	calendar := app.New(config.App.Name, logg, storage)

	server := internalhttp.NewServer(logg, calendar)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
