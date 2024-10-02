package dtos

import "time"

type Property struct {
	Type           string    `json:"@type"`
	IdLink         string    `json:"@id"`
	Id             string    `json:"id"`
	OrganizationID string    `json:"organization_id"`
	Slug           string    `json:"slug"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Properties struct {
	Type       string     `json:"@type"`
	ID         string     `json:"@id"`
	TotalItems int        `json:"total_items"`
	Items      []Property `json:"items"`
}
