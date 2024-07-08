// Package storage provides a simple in-memory storage for URL shortening.
package storage

import (
	"errors"

	urlgenerator "github.com/Mikeloangel/squasher/internal/urlGenerator"
)

// Storage represents an in-memory storage for shortened URLs.
type InMemoryStorage struct {
	data map[string]string
}

// NewStorage creates a new instance of In memory `Storage.
func NewInMemoryStorage() Storager {
	return &InMemoryStorage{
		data: make(map[string]string),
	}
}

// Set stores the given URL and returns a shortened version of it.
// If the URL is already stored, it returns the existing shortened version.
func (s *InMemoryStorage) StoreURL(url string) (StorageItem, error) {
	short := urlgenerator.HashURL(url)
	_, ok := s.data[short]
	if !ok {
		s.data[short] = url
	}

	return StorageItem{
		URL:     url,
		Shorten: short,
	}, nil
}

// Get retrieves the original URL for the given shortened version.
// It returns an error if the shortened URL is not found.
func (s *InMemoryStorage) FetchURL(short string) (StorageItem, error) {
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

// Init starts InMemoryStorage
func (s *InMemoryStorage) Init() error {
	return nil
}

// MultiStoreUrl creates for slice of links with options shorten version
func (s *InMemoryStorage) MultiStoreUrl(items *[]StorageItemOptionsInterface) error {
	for i, v := range *items {
		si, err := s.StoreURL(v.GetStorageItem().URL)
		if err != nil {
			return err
		}
		(*items)[i].GetStorageItem().Shorten = si.Shorten
	}
	return nil
}
