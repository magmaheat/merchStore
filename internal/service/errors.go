package service

import "errors"

var (
	ErrCannotParseToken = errors.New("cannot parse token")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrCannotSignToken  = errors.New("cannot sign token")
	ErrItemNotFound     = errors.New("item not found")
	ErrUserIdNotFound   = errors.New("user id not found")
	ErrBalanceTooLow    = errors.New("balance cannot be less than or equal to zero")
)
