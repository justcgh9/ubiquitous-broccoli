package storage

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrConflict       = errors.New("conflict")
	ErrNoPendingRequest = errors.New("no pending request")
)
