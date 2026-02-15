package grpc

import (
	"context"
	"github.com/Uranury/exploreMicro/service1/internal/storage"
	"github.com/Uranury/exploreMicro/service1/proto/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type userService struct {
	pb.UnimplementedUserServiceServer
	store storage.Store
}

func NewUserService(store storage.Store) pb.UserServiceServer {
	return &userService{
		store: store,
	}
}

func (s *userService) GetUser(_ context.Context, request *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, exists := s.store.Get(uint(request.UserId))
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &pb.GetUserResponse{
		Id:      uint32(user.ID),
		Name:    user.Name,
		Balance: user.Balance,
		Age:     uint32(user.Age),
	}, nil
}

func (s *userService) UpdateBalance(_ context.Context, request *pb.UpdateBalanceRequest) (*pb.UpdateBalanceResponse, error) {
	user, exists := s.store.Get(uint(request.UserId))
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	user.Balance = request.NewBalance
	s.store.Save(user)

	return &pb.UpdateBalanceResponse{
		Id:      uint32(user.ID),
		Name:    user.Name,
		Balance: user.Balance,
		Age:     uint32(user.Age),
	}, nil
}
