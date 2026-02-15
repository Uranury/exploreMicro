package service

import (
	"context"
	"github.com/Uranury/exploreMicro/service1/internal/models"
)

type User interface {
	CreateUser(ctx context.Context, name string, balance float64, age uint) (*models.User, error)
	GetUser(ctx context.Context, id uint) (*models.User, error)
	GetUsers(ctx context.Context) []*models.User
	UpdateUser(ctx context.Context, id uint, name *string, balance *float64, age *uint) (*models.User, error)
	UpdateBalance(ctx context.Context, id uint, balance float64) (*models.User, error)
}
