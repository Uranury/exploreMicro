package service

import (
	"context"
	"fmt"
	pb "github.com/Uranury/exploreMicro/service1/proto/pb"
	"github.com/Uranury/exploreMicro/service2/internal/models"
	"github.com/Uranury/exploreMicro/service2/internal/storage"
	"sync/atomic"
)

type Service interface {
	CreateOrder(ctx context.Context, userID uint, item string, price float64) (*models.Order, error)
	GetOrder(ctx context.Context, id uint) (*models.Order, error)
	ListOrders(ctx context.Context) ([]*models.Order, error)
	CancelOrder(ctx context.Context, id uint) error
}

type service struct {
	store  storage.Store
	client pb.UserServiceClient
	nextID atomic.Uint64
}

func NewService(store storage.Store, userClient pb.UserServiceClient) Service {
	return &service{
		store:  store,
		client: userClient,
	}
}

func (s *service) CreateOrder(ctx context.Context, userID uint, item string, price float64) (*models.Order, error) {
	id := uint(s.nextID.Add(1))
	userResp, err := s.client.GetUser(ctx, &pb.GetUserRequest{UserId: uint32(userID)})
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if userResp.Balance < price {
		return nil, fmt.Errorf("insufficient balance")
	}
	updatedUser, err := s.client.UpdateBalance(ctx, &pb.UpdateBalanceRequest{UserId: uint32(userID), NewBalance: userResp.Balance - price})
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	ord := models.Order{
		ID:     id,
		UserID: userID,
		Item:   item,
		Price:  price,
		User: &models.User{
			ID:      uint(updatedUser.Id),
			Name:    updatedUser.Name,
			Age:     uint(updatedUser.Age),
			Balance: updatedUser.Balance,
		},
	}
	s.store.Save(&ord)
	return &ord, nil
}

func (s *service) GetOrder(_ context.Context, id uint) (*models.Order, error) {
	if ord, exists := s.store.Get(id); exists {
		return ord, nil
	}
	return nil, fmt.Errorf("order with id %d not found", id)
}

func (s *service) ListOrders(_ context.Context) ([]*models.Order, error) {
	ords := s.store.List()
	return ords, nil
}

func (s *service) CancelOrder(ctx context.Context, id uint) error {
	ord, err := s.GetOrder(ctx, id)
	if err != nil {
		return err
	}
	user, err := s.client.GetUser(ctx, &pb.GetUserRequest{UserId: uint32(ord.UserID)})
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	_, err = s.client.UpdateBalance(ctx, &pb.UpdateBalanceRequest{UserId: uint32(ord.UserID), NewBalance: user.Balance + ord.Price})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	s.store.Delete(ord.ID)
	return nil
}
