package storage

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// StorageItemWithCorrelationID is a wrapper for StorageItem to utilize batch inserion api
// endpoint to pass correlation_id from request to response
type StorageItemWithCorrelationID struct {
	StorageItem
	CorrelationID string `json:"correlation_id"`
	options       []interface{}
}

// SetOptions adds CorrelationID param and options
func (m *StorageItemWithCorrelationID) SetOptions(options ...interface{}) error {
	m.options = options
	opt := options[0]
	switch v := opt.(type) {
	case int:
		m.CorrelationID = strconv.Itoa(v)
	case string:
		m.CorrelationID = v
	default:
		return fmt.Errorf("multiStoreStorageItem cant convert correlation_od")
	}
	return nil
}

// GetOptions return options
func (m *StorageItemWithCorrelationID) GetOptions() ([]interface{}, error) {
	return m.options, nil
}

// GetStorageItem provides access to StorageItem
func (m *StorageItemWithCorrelationID) GetStorageItem() *StorageItem {
	return &m.StorageItem
}

// MarshalJSON adds Marshaler interface support
func (m *StorageItemWithCorrelationID) MarshalJSON() ([]byte, error) {
	cid := m.CorrelationID
	shorten := m.Shorten

	return json.Marshal(&struct {
		Cid      string `json:"correlation_id"`
		Original string `json:"short_url"`
	}{
		Cid:      cid,
		Original: shorten,
	})
}
