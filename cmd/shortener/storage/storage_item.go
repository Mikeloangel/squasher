package storage

// StorageItem represents one item from storage
type (
	StorageItem struct {
		ID      int64  `json:"id"`
		URL     string `json:"original_url"`
		Shorten string `json:"short_url"`
	}
)
