package user

import "errors"

var (
	ErrPhoneExists  = errors.New("user with this phone already exists")
	ErrUserNotFound = errors.New("user not found")
)
