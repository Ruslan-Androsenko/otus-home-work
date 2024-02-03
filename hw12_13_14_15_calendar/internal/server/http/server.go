package internalhttp

import (
	"net/http"
	"time"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server"
)

var logg server.Logger

type Server struct {
	app    server.Application
	server http.Server
}

func NewServer(config server.Conf, app server.Application, logger server.Logger) *Server {
	address := config.GetHTTPAddress()
	logg = logger

	return &Server{
		app: app,
		server: http.Server{
			Addr:              address,
			ReadHeaderTimeout: time.Second * 3,
		},
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc(homePage, s.homePageHandler)
	mux.HandleFunc(helloPage, s.helloPageHandler)

	mux.HandleFunc(createEvent, s.createEventHandler)
	mux.HandleFunc(updateEvent, s.updateEventHandler)
	mux.HandleFunc(deleteEvent, s.deleteEventHandler)

	mux.HandleFunc(getEventByID, s.getEventHandler)
	mux.HandleFunc(getEventOfDay, s.getEventsOfDayHandler)
	mux.HandleFunc(getEventOfWeek, s.getEventsOfWeekHandler)
	mux.HandleFunc(getEventOfMonth, s.getEventsOfMonthHandler)

	s.server.Handler = loggingMiddleware(mux)
	return s.server.ListenAndServe()
}

func (s *Server) Stop() error {
	return s.server.Close()
}
