package internalhttp

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
)

type Server struct {
	app    Application
	server http.Server
}

type ServerConf struct {
	Host string
	Port int
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

type Application interface {
	UpdateEvent(ctx context.Context, id string, event storage.Event) error
	DeleteEvent(ctx context.Context, id string) error
	GetEventById(id string) (storage.Event, error)

	GetEventsOfDay(date time.Time) ([]storage.Event, error)
	GetEventsOfWeek(date time.Time) ([]storage.Event, error)
	GetEventsOfMonth(date time.Time) ([]storage.Event, error)
}

var logg Logger

func NewServer(logger Logger, app Application, config ServerConf) *Server {
	address := net.JoinHostPort(config.Host, strconv.Itoa(config.Port))
	logg = logger

	return &Server{
		app: app,
		server: http.Server{
			Addr: address,
		},
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", s.rootHandler)
	mux.HandleFunc("/hello", s.helloHandler)

	s.server.Handler = loggingMiddleware(mux)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Close()
}
