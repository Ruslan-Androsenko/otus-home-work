package grpc

import (
	"net"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/api/proto"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var logg server.Logger

type Server struct {
	proto.UnimplementedEventServiceServer

	app      server.Application
	grpc     *grpc.Server
	listener net.Listener
}

func NewServer(config server.Conf, app server.Application, logger server.Logger) *Server {
	logg = logger

	listener, err := net.Listen("tcp", config.GetGrpcAddress())
	if err != nil {
		logg.Fatalf("Failed to listen: %v", err)
	}

	return &Server{
		app:      app,
		grpc:     grpc.NewServer(),
		listener: listener,
	}
}

func (s *Server) Start() error {
	// Register service on gRPC server
	proto.RegisterEventServiceServer(s.grpc, s)

	// Register reflection service on gRPC server
	reflection.Register(s.grpc)

	return s.grpc.Serve(s.listener)
}

func (s *Server) Stop() error {
	s.grpc.Stop()
	return s.listener.Close()
}
