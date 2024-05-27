package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/api"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/config"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"
)

var (
	configFile string
	logg       *logger.Logger
)

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if hasVersionCommand() {
		printVersion()
		return
	}

	conf := config.NewConfig(configFile)
	logg = logger.New(conf.Logger.Level)
	storage := conf.GetStorage(logg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if conf.Storage.Type == config.StorageTypeDataBase {
		defer func() {
			if err := storage.Close(); err != nil {
				logg.Fatalf("Can't close database connection. Error: %v", err)
			}
		}()

		if err := storage.Connect(ctx); err != nil {
			logg.Fatalf("Can't connect to database. Error: %v", err)
		}

		if hasMigrationUpCommand() {
			config.MigrationUp(ctx, logg)
			return
		} else if hasMigrationDownCommand() {
			config.MigrationDown(ctx, logg)
			return
		}
	}

	calendar := app.New(logg, storage)
	server := api.NewServer(conf.Server, calendar, logg)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		server.Stop(ctx)
	}()

	logg.Info("calendar is running...")
	server.Start(ctx, cancel)
}
