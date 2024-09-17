package properties

type Properties struct {
	Type       string     `json:"@type"`
	ID         string     `json:"@id"`
	Links      Links      `json:"@links"`
	TotalItems int        `json:"total_items"`
	Items      []Property `json:"items"`
}
