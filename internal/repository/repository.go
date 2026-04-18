package repository

import (
	"PipisaBot/internal/model"
	"sync"
	"time"
)

type UserInMemory map[int64]map[int64]model.User

type UserRepository struct {
	m  UserInMemory
	mu sync.RWMutex
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		m:  make(UserInMemory),
		mu: sync.RWMutex{},
	}
}

func (r *UserRepository) CreateUser(user, chat int64) (model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	chatUsers, ok := r.m[chat]
	if !ok {
		chatUsers = make(map[int64]model.User)
		r.m[chat] = chatUsers
	}

	if _, ok := chatUsers[user]; ok {
		return model.User{}, AlreadyExists
	}

	u := model.User{}
	chatUsers[user] = u
	return u, nil
}

func (r *UserRepository) UpdateUser(user, chat, newLength int64, newTime time.Time) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	u, ok := r.m[chat][user]
	if !ok {
		return NotFound
	}

	u.Length = newLength
	u.LastUpdate = newTime

	r.m[chat][user] = u

	return nil
}

func (r *UserRepository) UpdateUserStatistic(user, chat, newLength int64) error {
	r.mu.Lock()

	defer r.mu.Unlock()

	u, ok := r.m[chat][user]
	if !ok {
		return NotFound
	}

	u.Length = newLength
	r.m[chat][user] = u

	return nil
}

func (r *UserRepository) GetUser(user, chat int64) (model.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.m[chat][user]
	if !ok {
		return model.User{}, NotFound
	}

	return u, nil
}

func (r *UserRepository) ListUsers(chat int64) []model.UserRecord {
	r.mu.RLock()
	defer r.mu.RUnlock()

	chatUsers := r.m[chat]
	users := make([]model.UserRecord, 0, len(chatUsers))
	for id, u := range chatUsers {
		users = append(users, model.UserRecord{
			ID:   id,
			User: u,
		})
	}

	return users
}
