package apperrors

import "errors"

var ErrItemAlreadyExists = errors.New("item already exists")

func IsErrItemAlreadyExists(err error) bool {
	return errors.Is(err, ErrItemAlreadyExists)
}

func NotErrItemAlreadyExists(err error) bool {
	return !errors.Is(err, ErrItemAlreadyExists)
}
