// Package storage provides a simple in-memory storage for URL shortening.
package storage

import (
	"errors"
	"fmt"
	"hash/fnv"
)

// Storage represents an in-memory storage for shortened URLs.
type Storage struct {
	data map[string]string
}

// NewStorage creates a new instance of Storage.
func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

// Set stores the given URL and returns a shortened version of it.
// If the URL is already stored, it returns the existing shortened version.
func (s *Storage) Set(url string) StorageItem {
	short := s.generateShortURL(url)
	_, ok := s.data[short]
	if !ok {
		s.data[short] = url
	}

	return StorageItem{
		URL:     url,
		Shorten: short,
	}
}

// Get retrieves the original URL for the given shortened version.
// It returns an error if the shortened URL is not found.
func (s *Storage) Get(short string) (StorageItem, error) {
	value, ok := s.data[short]
	if !ok {
		return StorageItem{}, errors.New("link not found")
	}

	si := StorageItem{
		URL:     value,
		Shorten: short,
	}
	return si, nil
}

// generateShortURL generates a shortened version of the given URL
// using the FNV-1a hash function.
func (s *Storage) generateShortURL(url string) string {
	hasher := fnv.New32a()
	hasher.Write([]byte(url))
	return fmt.Sprintf("%x", hasher.Sum32())
}
