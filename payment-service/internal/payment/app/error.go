package app

import "errors"

var (
	ErrEmailEmpty     = errors.New("email is required")
	ErrAmountInvalid  = errors.New("amount must be greater than zero")
	ErrStatusEmpty    = errors.New("status is required")
	ErrInsertFailed   = errors.New("failed to insert payment")
	ErrInvalidPayload = errors.New("invalid request format")
)
