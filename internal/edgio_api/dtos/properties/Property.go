package properties

import "time"

// Property represents a single property in the items list.
type Property struct {
	Type           string    `json:"@type"`
	ID             string    `json:"@id"`
	PropertyID     string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Slug           string    `json:"slug"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
