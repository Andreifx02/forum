package postrgres

import "errors"

var (
	UserNotFound = errors.New("User not found")
	UserAlreadyExists = errors.New("User already exists")
)