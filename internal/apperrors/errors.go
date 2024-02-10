package apperrors

import "errors"

var (
	EntityAlreadyExist = errors.New("entity already exist")
	IncorrectPassword  = errors.New("incorrect password")
	EntityNotFound     = errors.New("entity not found")
	InvalidToken       = errors.New("invalid token")
)
