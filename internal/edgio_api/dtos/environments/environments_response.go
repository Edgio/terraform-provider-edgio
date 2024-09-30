package environments

type EnvironmentsResponse struct {
	Type       string            `json:"@type"`
	Id         string            `json:"@id"`
	Links      EnvironmentsLinks `json:"@links"`
	TotalItems int               `json:"total_items"`
	Items      []Environment     `json:"items"`
}
