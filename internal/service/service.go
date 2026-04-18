package service

import (
	"PipisaBot/internal/repository"
	"crypto/rand"
	"errors"
	"math/big"
	"time"
)

type BoostResult struct {
	Delta     int64
	Length    int64
	Rank      int
	Penalized bool
}

type UserService struct {
	repository Repository
}

func NewUserService(repository Repository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) Boost(user, chat int64) (BoostResult, error) {
	now := time.Now()
	u, err := s.repository.GetUser(user, chat)
	if err != nil {
		if !errors.Is(err, repository.NotFound) {
			return BoostResult{}, err
		}

		u, err = s.repository.CreateUser(user, chat)
		if err != nil && !errors.Is(err, repository.AlreadyExists) {
			return BoostResult{}, err
		}
		if errors.Is(err, repository.AlreadyExists) {
			u, err = s.repository.GetUser(user, chat)
			if err != nil {
				return BoostResult{}, err
			}
		}
	}

	canBoost := u.LastUpdate.IsZero() || now.Sub(u.LastUpdate) >= 24*time.Hour

	var delta int64
	if canBoost {
		delta, err = randomBonus()
	} else {
		delta, err = randomAntiBonus()
	}
	if err != nil {
		return BoostResult{}, ErrGenerate
	}

	newLength := u.Length + delta
	if canBoost {
		err = s.repository.UpdateUser(user, chat, newLength, now)
	} else {
		err = s.repository.UpdateUserStatistic(user, chat, newLength)
	}
	if err != nil {
		return BoostResult{}, err
	}

	rank := 1
	users := s.repository.ListUsers(chat)
	for _, candidate := range users {
		if candidate.ID == user {
			continue
		}
		if candidate.User.Length > newLength {
			rank++
		}
	}

	return BoostResult{
		Delta:     delta,
		Length:    newLength,
		Rank:      rank,
		Penalized: !canBoost,
	}, nil
}

func randomAntiBonus() (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(11))
	if err != nil {
		return 0, err
	}
	return n.Int64() - 15, nil
}

func randomBonus() (int64, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(21))
	if err != nil {
		return 0, err
	}
	return n.Int64() - 10, nil
}
