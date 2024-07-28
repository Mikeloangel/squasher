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

	// MultiStoreURL sets multiple items to storage
	// usage of StorageItemOptionsInterface allows to process multiple items with additional data
	MultiStoreURL(items *[]StorageItemOptionsInterface) error
}

// StorageItemOptionsInterface defines methods for StorageItem wrappers
// to handle additional fields for specific implementations
// see StorageItemWithCorrelationID for example
type StorageItemOptionsInterface interface {
	// SetOptions extends StorageItem with additional fields

	SetOptions(options ...interface{}) error

	// GetOptions returns options for wrapper
	GetOptions() ([]interface{}, error)

	// GetStorageItem return pointer to StorageItem
	GetStorageItem() *StorageItem
}
