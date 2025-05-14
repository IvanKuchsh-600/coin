package entities

import "errors"

var (
	ErrInvalidParams  = errors.New("Invalid Params")
	ErrInternalServer = errors.New("Server Error")
	ErrGetFunc        = errors.New("Func Error")
)
