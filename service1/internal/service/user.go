package service

import (
	"context"
	"fmt"
	"github.com/Uranury/exploreMicro/service1/internal/models"
	"github.com/Uranury/exploreMicro/service1/internal/storage"
	"sync/atomic"
)

type user struct {
	store  storage.Store
	nextID atomic.Uint64
}

func NewUser(store storage.Store) User {
	return &user{
		store: store,
	}
}

func (s *user) CreateUser(_ context.Context, name string, balance float64, age uint) (*models.User, error) {
	id := uint(s.nextID.Add(1))
	newUser := &models.User{
		ID:      id,
		Name:    name,
		Balance: balance,
		Age:     age,
	}
	s.store.Save(newUser)
	return newUser, nil
}

func (s *user) GetUser(_ context.Context, id uint) (*models.User, error) {
	if user, exists := s.store.Get(id); exists {
		return user, nil
	}
	return nil, fmt.Errorf("user not found")
}

func (s *user) GetUsers(_ context.Context) []*models.User {
	return s.store.List()
}

func (s *user) UpdateUser(ctx context.Context, id uint, name *string, balance *float64, age *uint) (*models.User, error) {
	usr, err := s.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if name != nil {
		usr.Name = *name
	}
	if age != nil {
		usr.Age = *age
	}
	if balance != nil {
		usr.Balance = *balance
	}
	s.store.Save(usr)
	return usr, nil
}

func (s *user) UpdateBalance(ctx context.Context, id uint, balance float64) (*models.User, error) {
	return s.UpdateUser(ctx, id, nil, &balance, nil)
}
