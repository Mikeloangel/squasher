package storage

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Mikeloangel/squasher/internal/config"
	urlgenerator "github.com/Mikeloangel/squasher/internal/urlGenerator"
)

// Storage represents an file storage for shortened URLs.
type FileStorage struct {
	fileLocation string
	writer       *bufio.Writer
}

// NewFileStorage creates a new instance of file storage
func NewFileStorage(cfg *config.Config) Storager {
	return &FileStorage{
		fileLocation: cfg.StorageFileLocation,
	}
}

// Set stores the given URL and returns a shortened version of it.
// If the URL is already stored, it returns the existing shortened version.
func (s *FileStorage) StoreURL(url string) (StorageItem, error) {
	short := urlgenerator.HashURL(url)

	// Tries to get already existed url
	si, err := s.FetchURL(short)
	if err == nil {
		return si, NewItemAlreadyExistsError(nil, url)
	}

	si = StorageItem{
		URL:     url,
		Shorten: short,
	}

	// Appends to file
	data, err := json.Marshal(&si)
	if err != nil {
		return si, err
	}

	if s.writer == nil {
		return si, errors.New("writer is not initilised")
	}

	if _, err = s.writer.Write(data); err != nil {
		return si, err
	}

	if err = s.writer.WriteByte('\n'); err != nil {
		return si, err
	}

	return si, s.writer.Flush()
}

// Get retrieves the original URL for the given shortened version.
// It returns an error if the shortened URL is not found.
func (s *FileStorage) FetchURL(short string) (StorageItem, error) {
	si := StorageItem{}

	file, err := os.Open(s.fileLocation)
	if err != nil {
		return si, fmt.Errorf("filestorage failed to open file at FetchURL")
	}
	defer file.Close()

	// Scanns for url
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		err := json.Unmarshal(scanner.Bytes(), &si)
		if err != nil {
			return si, err
		}
		if si.Shorten == short {
			return si, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return si, err
	}

	return si, errors.New("shortened URL not found")
}

// Init opens file to write
func (s *FileStorage) Init() error {
	file, err := os.OpenFile(s.fileLocation, os.O_WRONLY|os.O_APPEND|os.O_CREATE|os.O_SYNC, 0666)
	if err != nil {
		return err
	}

	s.writer = bufio.NewWriter(file)

	return nil
}

// MultiStoreURL creates for slice of links with options shorten version
func (s *FileStorage) MultiStoreURL(items *[]StorageItemOptionsInterface) error {
	for i, v := range *items {
		si, err := s.StoreURL(v.GetStorageItem().URL)
		var ae *ItemAlreadyExistsError
		if err != nil && !errors.As(err, &ae) {
			return err
		}
		(*items)[i].GetStorageItem().Shorten = si.Shorten
	}
	return nil
}
