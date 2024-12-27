package constants

import (
	"errors"
)

var (
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrEditConflict   = errors.New("edit conflict")
	ErrRecordNotFound = errors.New("record not found")
)
