package messages

import "errors"

const (
	MsgHealthOK   = "Health ok"
	WrongResponse = "Wrong"
)

var (
	ErrFailedUsername = errors.New("Failed finding and creating user")
	ErrNoToken        = errors.New("No token found")
)
