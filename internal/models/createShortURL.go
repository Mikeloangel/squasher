package models

import "github.com/asaskevich/govalidator"

// CreateShortURLRequest represents the request payload for creating a short URL.
// The URL field is required and must be a valid URL.
type CreateShortURLRequest struct {
	URL string `json:"url" valid:"required,url"`
}

// CreateShortURLResponse represents the response payload after creating a short URL.
// The Result field contains the shortened URL.
type CreateShortURLResponse struct {
	Result string `json:"result"`
}

// ValidateCreateShortURLRequest validates the CreateShortURLRequest struct.
// It ensures that the URL field is present and is a valid URL.
// Returns an error if validation fails.
func ValidateCreateShortURLRequest(req *CreateShortURLRequest) error {
	// Validate the struct based on the `valid` tags.
	_, err := govalidator.ValidateStruct(req)
	return err
}
