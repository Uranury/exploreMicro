package storage

import (
	"github.com/Uranury/exploreMicro/service2/internal/models"
	"sync"
)

type Store interface {
	Get(id uint) (*models.Order, bool)
	List() []*models.Order
	Save(order *models.Order)
	Delete(id uint)
}

type store struct {
	mu     sync.RWMutex
	orders map[uint]*models.Order
}

func NewStore() Store {
	return &store{
		orders: make(map[uint]*models.Order),
	}
}

func (s *store) Get(id uint) (*models.Order, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[id]
	return order, ok
}

func (s *store) List() []*models.Order {
	s.mu.RLock()
	defer s.mu.RUnlock()
	orders := make([]*models.Order, 0, len(s.orders))
	for _, user := range s.orders {
		orders = append(orders, user)
	}
	return orders
}

func (s *store) Save(order *models.Order) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.orders[order.ID] = order
}

func (s *store) Delete(id uint) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.orders, id)
}
