package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"

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

	conf := config.NewConfig(configFile)
	logg = logger.New(conf.Logger.Level)
	eventStorage := conf.GetStorage(logg)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	if conf.Storage.Type == config.StorageTypeDataBase {
		defer func() {
			if err := eventStorage.Close(); err != nil {
				logg.Fatalf("Can't close database connection. Error: %v", err)
			}
		}()

		if err := eventStorage.Connect(ctx); err != nil {
			logg.Fatalf("Can't connect to database. Error: %v", err)
		}
	}

	queue := conf.GetQueue()

	defer func() {
		if err := queue.Close(); err != nil {
			logg.Fatalf("Can't close Queue connection. Error: %v", err)
		}
	}()

	logg.Infof("Dialing sheduler on queue...")
	if err := queue.Connect(); err != nil {
		logg.Fatalf("Can't connect to Queue server. Error: %v", err)
	}

	go func() {
		logg.Errorf("closing on sheduler: %s", <-queue.NotifyClose())
	}()

	go func() {
		<-ctx.Done()
	}()

	notifications, err := eventStorage.GetEventsNotifications()
	if err != nil {
		logg.Errorf("cannot getting Events notifications. Error: %v", err)
	}

	for _, notification := range notifications {
		if err = queue.Publish(notification); err != nil {
			logg.Errorf("cannot publish message: %v. Error: %v", notification, err)
		}
	}
}
