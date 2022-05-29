package storage

import "errors"

var (
	ErrAlreadyExist = errors.New("already exist")
)

func IsErrAlreadyExist(err error) bool {
	return errors.Is(err, ErrAlreadyExist)
}
