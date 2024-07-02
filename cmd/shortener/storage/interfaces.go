// Package interfaces defines interfaces for the URL shortener service.
package storage

// Storager defines the methods required for a URL storage system.
type Storager interface {
	// Get retrieves the original URL for the given shortened version.
	// It returns an error if the shortened URL is not found.
	FetchURL(short string) (StorageItem, error)

	// Set stores the given URL and returns a shortened version of it.
	// If the URL is already stored, it returns the existing shortened version.
	StoreURL(url string) (StorageItem, error)

	// Init initializes storage
	Init() error
}
