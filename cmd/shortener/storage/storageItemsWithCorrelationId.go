package storage

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// StorageItemWithCorrelationId is a wrapper for StorageItem to utilize batch inserion api
// endpoint to pass correlation_id from request to response
type StorageItemWithCorrelationId struct {
	StorageItem
	CorrelationId string `json:"correlation_id"`
	options       []interface{}
}

// SetOptions adds CorrelationId param and options
func (m *StorageItemWithCorrelationId) SetOptions(options ...interface{}) error {
	m.options = options
	opt := options[0]
	switch v := opt.(type) {
	case int:
		m.CorrelationId = strconv.Itoa(v)
	case string:
		m.CorrelationId = v
	default:
		return fmt.Errorf("multiStoreStorageItem cant convert correlation_od")
	}
	return nil
}

// GetOptions return options
func (m *StorageItemWithCorrelationId) GetOptions() ([]interface{}, error) {
	return m.options, nil
}

// GetStorageItem provides access to StorageItem
func (m *StorageItemWithCorrelationId) GetStorageItem() *StorageItem {
	return &m.StorageItem
}

// MarshalJSON adds Marshaler interface support
func (m *StorageItemWithCorrelationId) MarshalJSON() ([]byte, error) {
	cid := m.CorrelationId
	shorten := m.Shorten

	return json.Marshal(&struct {
		Cid      string `json:"correlation_id"`
		Original string `json:"short_url"`
	}{
		Cid:      cid,
		Original: shorten,
	})
}
