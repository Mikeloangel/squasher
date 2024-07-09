package storage

import "fmt"

// ItemAlreadyExistsError is an error wrapper for when item already in storage
type ItemAlreadyExistsError struct {
	OriginalURL string
	Err         error
}

// Error implements Error interface
func (conflict ItemAlreadyExistsError) Error() string {
	return fmt.Sprintf("already exists: %s", conflict.OriginalURL)
}

// Unwrap implements Error interface
func (conflict ItemAlreadyExistsError) Unwrap() error {
	return conflict.Err
}

// NewItemAlreadyExistsError returns new ItemAlreadyExistsError error
func NewItemAlreadyExistsError(err error, originalURL string) *ItemAlreadyExistsError {
	return &ItemAlreadyExistsError{
		OriginalURL: originalURL,
		Err:         err,
	}
}
