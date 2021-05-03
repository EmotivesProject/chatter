package messages

import "errors"

const WrongResponse = "Wrong"

var (
	ErrFailedUsername = errors.New("Failed finding and creating user")
	ErrFailedToType   = errors.New("Failed parsing to a type")
	ErrNoToken        = errors.New("No token found")
)
