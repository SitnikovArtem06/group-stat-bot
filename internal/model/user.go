package model

import "time"

type User struct {
	Length     int64
	LastUpdate time.Time
}

type UserRecord struct {
	ID   int64
	User User
}
