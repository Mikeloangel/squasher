package storage

import (
	"bufio"
	"encoding/json"
	"errors"
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
func (s *FileStorage) Set(url string) (StorageItem, error) {
	short := urlgenerator.GenerateShortURL(url)

	// Tries to get already existed file
	si, err := s.Get(short)
	if err == nil {
		return si, nil
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
func (s *FileStorage) Get(short string) (StorageItem, error) {
	si := StorageItem{}

	file, err := os.Open(s.fileLocation)
	if err != nil {
		return si, err
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
