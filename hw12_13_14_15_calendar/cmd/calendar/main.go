package main

import (
	"context"
	"database/sql"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/api"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"
)

var (
	configFile string
	dbConn     *sql.DB
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

	config := NewConfig()
	logg = logger.New(config.Logger.Level)
	storage := config.GetStorage()

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if config.Storage.Type == StorageTypeDataBase {
		defer func() {
			err := storage.Close()
			if err != nil {
				logg.Fatalf("Can't close database connection. Error: %v", err)
			}
		}()

		err := storage.Connect(ctx)
		if err != nil {
			logg.Fatalf("Can't connect to database. Error: %v", err)
		}

		if hasMigrationUpCommand() {
			migrationUp(ctx)
			return
		} else if hasMigrationDownCommand() {
			migrationDown(ctx)
			return
		}
	}

	calendar := app.New(logg, storage)
	server := api.NewServer(config.Server, calendar, logg)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		server.Stop(ctx)
	}()

	logg.Info("calendar is running...")
	server.Start(ctx, cancel)
}
