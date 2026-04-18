package service

import (
	"PipisaBot/internal/model"
	"time"
)

type Repository interface {
	CreateUser(user, chat int64) (model.User, error)
	UpdateUser(user, chat, newLength int64, newTime time.Time) error
	GetUser(user, chat int64) (model.User, error)
	UpdateUserStatistic(user, chat, newLength int64) error
	ListUsers(chat int64) []model.UserRecord
}
