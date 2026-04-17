package repository

import (
	"PipisaBot/internal/model"
	"sync"
	"time"
)

type UserInMemory map[int64]model.User

type UserRepository struct {
	m  UserInMemory
	mu sync.RWMutex
}

func NewUserRepository() *UserRepository {

	Um := make(map[int64]model.User)

	return &UserRepository{
		m:  Um,
		mu: sync.RWMutex{},
	}

}

func (r *UserRepository) CreateUser(user int64) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.m[user]; ok {
		return model.User{}, AlreadyExists
	}

	u := model.User{}
	r.m[user] = u
	return u, nil
}

func (r *UserRepository) UpdateUser(user int64, newLength int64, newTime time.Time) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	u, ok := r.m[user]
	if !ok {
		return NotFound
	}

	u.Length = newLength
	u.LastUpdate = newTime

	r.m[user] = u

	return nil
}

func (r *UserRepository) UpdateUserStatistic(user int64, newLength int64) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	u, ok := r.m[user]
	if !ok {
		return NotFound
	}

	u.Length = newLength
	r.m[user] = u

	return nil
}

func (r *UserRepository) GetUser(user int64) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.m[user]
	if !ok {
		return model.User{}, NotFound
	}

	return u, nil
}

func (r *UserRepository) ListUsers() []model.UserRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]model.UserRecord, 0, len(r.m))
	for id, u := range r.m {
		users = append(users, model.UserRecord{
			ID:   id,
			User: u,
		})
	}

	return users
}
