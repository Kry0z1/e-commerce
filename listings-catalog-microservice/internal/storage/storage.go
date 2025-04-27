package storage

import "errors"

var (
	ErrListingNotFound = errors.New("listing with such id not found")
	ErrUserNotFound    = errors.New("user with such id not found")
)
