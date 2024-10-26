package handlers

import (
	"context"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/SversusN/keeper/pkg/grpc"
)

// Ping – метод для проверки работоспособности сервера.
func (s *Server) Ping(ctx context.Context, in *pb.PingRequest) (*pb.PingResponse, error) {
	err := s.Storage.HealthCheck()
	if err != nil {
		return nil, status.Errorf(codes.Internal, http.StatusText(http.StatusInternalServerError))
	}

	// time.Sleep(5 * time.Second) //nolint:gomnd // for grasfull shutdown test
	return &pb.PingResponse{
		Result: "OK",
	}, nil
}
