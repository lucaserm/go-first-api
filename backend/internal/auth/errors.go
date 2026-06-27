package auth

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailConflict      = errors.New("this email is in use")
	ErrUsernameConflict   = errors.New("this username is in use")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
