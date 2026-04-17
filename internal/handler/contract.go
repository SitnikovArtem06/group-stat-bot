package handler

import "PipisaBot/internal/service"

type Service interface {
	Boost(user int64) (service.BoostResult, error)
}
