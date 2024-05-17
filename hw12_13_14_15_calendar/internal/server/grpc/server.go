package grpc

import (
	"context"
	"fmt"
	"net"

	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/api/proto"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/server"
	"github.com/Ruslan-Androsenko/otus-home-work/hw12_13_14_15_calendar/internal/storage"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
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

	// Настраиваем интерцептор логирования запросов
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(logging.PayloadSent),
	}
	opts := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), loggingOpts...),
		),
	}

	return &Server{
		app:      app,
		grpc:     grpc.NewServer(opts...),
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

func (s *Server) CreateEvent(ctx context.Context, req *proto.CreateEventRequest) (*proto.CreateEventResponse, error) {
	response := &proto.CreateEventResponse{}
	err := s.app.CreateEvent(ctx, storage.MakeStorageEvent(req.Event))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to create new event: %v, Error: %v", req.Event, err)
		response.Result = errorMessage
		logg.Error(errorMessage)
	} else {
		response.Result = "Ok"
	}

	return response, err
}

func (s *Server) UpdateEvent(ctx context.Context, req *proto.UpdateEventRequest) (*proto.UpdateEventResponse, error) {
	response := &proto.UpdateEventResponse{}
	err := s.app.UpdateEvent(ctx, req.Id, storage.MakeStorageEvent(req.Event))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to update eventId: %s, event: %v, Error: %v", req.Id, req.Event, err)
		response.Result = errorMessage
		logg.Error(errorMessage)
	} else {
		response.Result = "Ok"
	}

	return response, err
}

func (s *Server) DeleteEvent(ctx context.Context, req *proto.GetEventRequest) (*proto.DeleteEventResponse, error) {
	response := &proto.DeleteEventResponse{}
	err := s.app.DeleteEvent(ctx, req.Id)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to delete eventId: %s, Error: %v", req.Id, err)
		response.Result = errorMessage
		logg.Error(errorMessage)
	} else {
		response.Result = "Ok"
	}

	return response, err
}

func (s *Server) GetEventByID(_ context.Context, req *proto.GetEventRequest) (*proto.GetEventResponse, error) {
	response := &proto.GetEventResponse{}
	event, err := s.app.GetEventByID(req.Id)
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get eventId: %s, Error: %v", req.Id, err)
		response.Result = errorMessage
		logg.Error(errorMessage)
	} else {
		response.Event = storage.MakeProtoEvent(event)
		response.Result = "Ok"
	}

	return response, err
}

func (s *Server) GetEventsOfDay(_ context.Context, req *proto.GetEventsListRequest) (*proto.GetEventsListResponse, error) {
	response := &proto.GetEventsListResponse{}
	events, err := s.app.GetEventsOfDay(storage.MakeTime(req.Date))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get events of day by date: %v, Error: %v", req.Date, err)
		response.Result = errorMessage
		logg.Error(errorMessage)
	} else {
		response.Events = storage.MakeProtoEventsList(events)
		response.Result = "Ok"
	}

	return response, err
}

func (s *Server) GetEventsOfWeek(_ context.Context, req *proto.GetEventsListRequest) (*proto.GetEventsListResponse, error) {
	response := &proto.GetEventsListResponse{}
	events, err := s.app.GetEventsOfWeek(storage.MakeTime(req.Date))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get events of week by date: %v, Error: %v", req.Date, err)
		response.Result = errorMessage
		logg.Error(errorMessage)
	} else {
		response.Events = storage.MakeProtoEventsList(events)
		response.Result = "Ok"
	}

	return response, err
}

func (s *Server) GetEventsOfMonth(_ context.Context, req *proto.GetEventsListRequest) (*proto.GetEventsListResponse, error) {
	response := &proto.GetEventsListResponse{}
	events, err := s.app.GetEventsOfMonth(storage.MakeTime(req.Date))
	if err != nil {
		errorMessage := fmt.Sprintf("Failed to get events of month by date: %v, Error: %v", req.Date, err)
		response.Result = errorMessage
		logg.Error(errorMessage)
	} else {
		response.Events = storage.MakeProtoEventsList(events)
		response.Result = "Ok"
	}

	return response, err
}
