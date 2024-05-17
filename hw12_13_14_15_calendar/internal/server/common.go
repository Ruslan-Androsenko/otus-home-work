package server

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Conf struct {
	Host     string `toml:"host"`
	GrpcPort int    `toml:"grpc_port"`
	HTTPPort int    `toml:"http_port"`
}

func (s Conf) GetGrpcAddress() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.GrpcPort))
}

func (s Conf) GetHTTPAddress() string {
	return net.JoinHostPort(s.Host, strconv.Itoa(s.HTTPPort))
}

type Application interface {
	CreateEvent(ctx context.Context, event storage.Event) error
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error

	GetEventByID(id string) (storage.Event, error)
	GetEventsOfDay(date time.Time) ([]storage.Event, error)
	GetEventsOfWeek(date time.Time) ([]storage.Event, error)
	GetEventsOfMonth(date time.Time) ([]storage.Event, error)
}

type Logger interface {
	Fatal(msg string)
	Error(msg string)
	Warning(msg string)
	Info(msg string)
	Debug(msg string)

	Fatalf(format string, values ...any)
	Errorf(format string, values ...any)
	Warningf(format string, values ...any)
	Infof(format string, values ...any)
	Debugf(format string, values ...any)
}
