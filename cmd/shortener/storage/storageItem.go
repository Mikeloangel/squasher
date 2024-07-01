package storage

// StorageItem represents one item from storage
type StorageItem struct {
	URL     string `json:"original_url"`
	Shorten string `json:"short_url"`
}
