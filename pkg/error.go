package pkg

import "errors"

var (
	ErrHandleTypeNotExist  = errors.New("handle type error")
	CommandLenError        = errors.New("command len error")
	ErrInterfaceType       = errors.New("interface type error")
	ErrKeyLenZero          = errors.New("key len is zero")
	ErrCommandTypeNotExist = errors.New("command type not exist")
	ErrInvalidReadResponse = errors.New("invalid read response")
)
