package storage

import "errors"

var ErrNotFound = errors.New("not found")

type ErrNotFound2 string

func (err ErrNotFound2) Error() string {
	return string(err)
}

func IsErrNotFound(err error) bool {
	var nf2 ErrNotFound2
	return errors.Is(err, ErrNotFound) || errors.As(err, &nf2)
}

type ErrAlreadyExists string

func NewErrAlreadyExists(msg string) ErrAlreadyExists {
	return ErrAlreadyExists(msg)
}

func (err ErrAlreadyExists) Error() string {
	return string(err)
}
