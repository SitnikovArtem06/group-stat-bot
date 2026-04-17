package service

import (
	"PipisaBot/internal/model"
	"time"
)

type Repository interface {
	CreateUser(user int64) (model.User, error)
	UpdateUser(user int64, newLength int64, newTime time.Time) error
	GetUser(user int64) (model.User, error)
	UpdateUserStatistic(user int64, newLength int64) error
	ListUsers() []model.UserRecord
}
