package grpc

import (
	"context"
	"github.com/Uranury/exploreMicro/service1/internal/service"
	"github.com/Uranury/exploreMicro/service1/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userHandler struct {
	pb.UnimplementedUserServiceServer
	svc service.User
}

func NewUserService(svc service.User) pb.UserServiceServer {
	return &userHandler{
		svc: svc,
	}
}

func (s *userHandler) GetUser(ctx context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.svc.GetUser(ctx, uint(request.UserId))
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		Id:      uint32(user.ID),
		Name:    user.Name,
		Balance: user.Balance,
		Age:     uint32(user.Age),
	}, nil
}

func (s *userHandler) UpdateBalance(ctx context.Context, request *pb.UpdateBalanceRequest) (*pb.UpdateBalanceResponse, error) {
	user, err := s.svc.UpdateBalance(ctx, uint(request.UserId), request.NewBalance)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update balance")
	}
	return &pb.UpdateBalanceResponse{
		Id:      uint32(user.ID),
		Name:    user.Name,
		Balance: user.Balance,
		Age:     uint32(user.Age),
	}, nil
}
