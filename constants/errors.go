package constants

import "errors"

var (
	ErrInvalidAuth  = errors.New("username or password is wrong")
	ErrEmailExist   = errors.New("email address already exist")
)