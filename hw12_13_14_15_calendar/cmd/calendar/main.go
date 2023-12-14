package main

import (
	"context"
	"database/sql"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/app"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"

	internalhttp "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server/http"
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
	server := internalhttp.NewServer(logg, calendar, config.Server)

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
