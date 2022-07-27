package errutil

import "errors"

var (
	ErrUnKnown          = errors.New("unknows error")
	ErrAccountExists    = errors.New("account exists")
	ErrAuthFailed       = errors.New("auth failed")
	ErrIllegalLoginType = errors.New("illegal login type")
	ErrNotFound         = errors.New("not found")
	ErrWrongType        = errors.New("wrong type")
	ErrWrongPassword    = errors.New("wrong account or password")
	ErrInitFailed       = errors.New("initialize failed")
	ErrNotImplemented   = errors.New("not implemented")
)
