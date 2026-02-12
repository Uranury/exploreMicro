package service

import (
	"context"
	"fmt"
	"github.com/Uranury/exploreMicro/service2/internal/http_pack"
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
	client http_pack.UserClient
	nextID atomic.Uint64
}

func NewService(store storage.Store, client http_pack.UserClient) Service {
	return &service{
		store:  store,
		client: client,
	}
}

func (s *service) CreateOrder(ctx context.Context, userID uint, item string, price float64) (*models.Order, error) {
	id := uint(s.nextID.Add(1))
	user, err := s.client.Get(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	if user.Balance < price {
		return nil, fmt.Errorf("insufficient balance")
	}
	user, err = s.client.Patch(ctx, userID, user.Balance-price)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}
	ord := models.Order{
		ID:     id,
		UserID: userID,
		Item:   item,
		Price:  price,
		User:   user,
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
	user, err := s.client.Get(ctx, ord.UserID)
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}
	_, err = s.client.Patch(ctx, ord.UserID, user.Balance+ord.Price)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	s.store.Delete(ord.ID)
	return nil
}
