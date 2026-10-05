package session

import "errors"

var (
	ErrInvalidSessionID = errors.New("invalid session id")
	ErrSessionNotFound  = errors.New("session not found")
	ErrUnknownRole      = errors.New("unknown role")
)
