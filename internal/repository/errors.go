package repository

import "errors"

var (
	NotFound      = errors.New("user not found")
	AlreadyExists = errors.New("user already exists")
)
