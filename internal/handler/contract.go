package handler

import "PipisaBot/internal/service"

type Service interface {
	Boost(user, chat int64) (service.BoostResult, error)
}
