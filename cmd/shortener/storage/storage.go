package storage

import (
	"errors"
	"fmt"
	"hash/fnv"
)

type Storage struct {
	data map[string]string
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[string]string),
	}
}

func (s *Storage) Set(url string) string {
	short := s.generateShortURL(url)
	_, ok := s.data[short]
	if !ok {
		s.data[short] = url
	}

	return short
}

func (s *Storage) Get(short string) (string, error) {
	value, ok := s.data[short]

	if !ok {
		return "", errors.New("link not found")
	}
	return value, nil
}

func (s *Storage) generateShortURL(url string) string {
	hasher := fnv.New32a()
	hasher.Write([]byte(url))
	return fmt.Sprintf("%x", hasher.Sum32())
}
