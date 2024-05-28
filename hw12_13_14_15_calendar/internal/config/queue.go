package config

import (
	"fmt"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/logger"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/queue/rabbitmq"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Connect() error
	Close() error
	NotifyClose() chan *amqp.Error
	Publish(message storage.Notification) error
	Consume() error
	Handle(done chan error, logg *logger.Logger)
}

type QueueConf struct {
	Host     string `toml:"queue_host"`
	Port     int    `toml:"queue_port"`
	User     string `toml:"queue_user"`
	Pass     string `toml:"queue_pass"`
	Exchange rabbitmq.ExchangeConf
	Period   PeriodConf
}

type PeriodConf struct {
	LaunchFrequency string `toml:"launch_frequency"`
}

// Получить строку соединения к серверу очередей.
func (config QueueConf) getDataSource() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/",
		config.User, config.Pass, config.Host, config.Port)
}
