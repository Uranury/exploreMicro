package storage

import (
	"github.com/Uranury/exploreMicro/service1/internal/models"
	"sync"
)

type Store interface {
	Get(id uint) (*models.User, bool)
	List() []*models.User
	Save(user *models.User)
}

type store struct {
	mu    sync.RWMutex
	users map[uint]*models.User
}

func NewStore() Store {
	return &store{
		users: make(map[uint]*models.User),
	}
}

func (s *store) Get(id uint) (*models.User, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	return user, ok
}

func (s *store) List() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	users := make([]*models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users
}

func (s *store) Save(user *models.User) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.users[user.ID] = user
}
