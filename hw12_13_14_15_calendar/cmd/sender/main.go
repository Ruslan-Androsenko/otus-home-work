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

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	queue := conf.GetQueue()

	defer func() {
		if err := queue.Close(); err != nil {
			logg.Fatalf("Can't close Queue connection. Error: %v", err)
		}
	}()

	logg.Infof("Dialing sender on queue...")
	if err := queue.Connect(); err != nil {
		logg.Fatalf("Can't connect to Queue server. Error: %v", err)
	}

	go func() {
		logg.Errorf("closing on sender: %s", <-queue.NotifyClose())
	}()

	go func() {
		<-ctx.Done()
	}()

	if err := queue.Consume(); err != nil {
		logg.Fatalf("Can't consume from Queue. Error: %v", err)
	}

	done := make(chan error)
	queue.Handle(done, logg)

	if err := <-done; err != nil {
		logg.Fatalf("Can't handle message from queue. Error: %v", err)
	}
}
