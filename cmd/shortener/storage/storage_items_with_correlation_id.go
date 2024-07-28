package storage

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/asaskevich/govalidator"
)

// StorageItemWithCorrelationID is a wrapper for StorageItem to utilize batch inserion api
// endpoint to pass correlation_id from request to response
type StorageItemWithCorrelationID struct {
	StorageItem
	CorrelationID string `json:"correlation_id" valid:"required"`
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

// ValidateStorageItemWithCorrelationIDRequest validates one item of StorageItemWithCorrelationID
func ValidateStorageItemWithCorrelationIDRequest(req *StorageItemWithCorrelationID) error {
	var err error
	if strings.TrimSpace(req.URL) == "" {
		return fmt.Errorf("original_url in request is required in each row for ID:" + req.CorrelationID)
	}

	isURL := govalidator.IsURL(req.URL)
	if !isURL {
		return fmt.Errorf("original_url in each row has to be valid url for ID:" + req.CorrelationID)
	}
	_, err = govalidator.ValidateStruct(req)
	return err
}

// ValidateStorageItemsWithCorrelationIDRequest validates array of StorageItemWithCorrelationID
func ValidateStorageItemsWithCorrelationIDRequest(items *[]StorageItemWithCorrelationID) error {
	for _, v := range *items {
		err := ValidateStorageItemWithCorrelationIDRequest(&v)
		if err != nil {
			return err
		}
	}
	return nil
}
