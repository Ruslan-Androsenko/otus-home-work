package api

import (
	"context"
	"os"
	"sync"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server"
	internalgrpc "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server/http"
)

type Server struct {
	grpc *internalgrpc.Server
	http *internalhttp.Server
}

var logg server.Logger

func NewServer(config server.Conf, app server.Application, logger server.Logger) *Server {
	logg = logger

	return &Server{
		grpc: internalgrpc.NewServer(config, app, logger),
		http: internalhttp.NewServer(config, app, logger),
	}
}

func (s *Server) Start(_ context.Context, cancel context.CancelFunc) {
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		logg.Info("grpc server is running...")
		if err := s.grpc.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	go func() {
		defer wg.Done()

		logg.Info("http server is running...")
		if err := s.http.Start(); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	}()

	wg.Wait()
}

func (s *Server) Stop(_ context.Context) {
	logg.Info("grpc server is stopped...")
	if err := s.grpc.Stop(); err != nil {
		logg.Error("failed to stop grpc server: " + err.Error())
	}

	logg.Info("http server is stopped...")
	if err := s.http.Stop(); err != nil {
		logg.Error("failed to stop http server: " + err.Error())
	}
}
