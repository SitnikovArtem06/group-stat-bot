package service

import "errors"

var (
	ErrNotFoundUser = errors.New("user not found")
	ErrGenerate     = errors.New("generate error bonus")
)
