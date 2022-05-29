package storage

import "errors"

var (
	ErrAlreadyExist = errors.New("already exist")
	ErrInvalidInfo  = errors.New("invalid info")
	ErrNotFound     = errors.New("not found")
)

func IsErrAlreadyExist(err error) bool {
	return errors.Is(err, ErrAlreadyExist)
}

func IsErrInvalidInfo(err error) bool {
	return errors.Is(err, ErrInvalidInfo)
}

func IsErrNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}
