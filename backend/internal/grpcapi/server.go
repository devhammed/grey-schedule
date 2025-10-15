package grpcapi

import (
	"context"
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	"github.com/devhammed/grey-schedule/backend/internal/models"
	"github.com/devhammed/grey-schedule/backend/internal/store"
	"github.com/devhammed/grey-schedule/backend/proto"
)

type Server struct {
	schedulepb.UnimplementedAppointmentServiceServer
	st     store.Store
	server *grpc.Server
}

func NewServer(s store.Store) *Server {
	gs := grpc.NewServer()

	srv := &Server{
		st:     s,
		server: gs,
	}

	schedulepb.RegisterAppointmentServiceServer(gs, srv)

	reflection.Register(gs)

	return srv
}

func (s *Server) CreateAppointment(_ context.Context, req *schedulepb.CreateAppointmentRequest) (*schedulepb.CreateAppointmentResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is required")
	}

	if req.GetTitle() == "" {
		return nil, status.Error(codes.InvalidArgument, models.ErrInvalidTitle.Error())
	}

	start, err := time.Parse(time.RFC3339, req.GetStart())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "start must be RFC3339 timestamp")
	}

	end, err := time.Parse(time.RFC3339, req.GetEnd())

	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "end must be RFC3339 timestamp")
	}

	appointment, err := s.st.Create(req.GetTitle(), start, end)

	if err != nil {
		switch {
		case errors.Is(err, store.ErrConflict):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		case errors.Is(err, models.ErrInvalidTitle), errors.Is(err, models.ErrInvalidTimeRange):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, "internal error")
		}
	}

	return &schedulepb.CreateAppointmentResponse{Appointment: toPb(appointment)}, nil
}

func (s *Server) ListAppointments(_ context.Context, _ *schedulepb.ListAppointmentsRequest) (*schedulepb.ListAppointmentsResponse, error) {
	list := s.st.List()
	out := make([]*schedulepb.Appointment, 0, len(list))

	for _, a := range list {
		out = append(out, toPb(a))
	}

	return &schedulepb.ListAppointmentsResponse{Appointments: out}, nil
}

func (s *Server) DeleteAppointment(_ context.Context, req *schedulepb.DeleteAppointmentRequest) (*schedulepb.DeleteAppointmentResponse, error) {
	if req == nil || req.GetId() == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	if err := s.st.Delete(req.GetId()); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &schedulepb.DeleteAppointmentResponse{}, nil
}

func (s *Server) Start(addr string) error {
	lis, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	return s.server.Serve(lis)
}

func toPb(a models.Appointment) *schedulepb.Appointment {
	return &schedulepb.Appointment{
		Id:        a.ID,
		Title:     a.Title,
		Start:     a.Start.Format(time.RFC3339),
		End:       a.End.Format(time.RFC3339),
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
	}
}
