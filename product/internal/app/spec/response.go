package spec

// ListResponse ...
type ListResponse struct {
	Products []GetResponse
}

// GetResponse ...
type GetResponse struct {
	ID                int     `json:"id"`
	Name              string  `json:"name"`
	Price             float64 `json:"price"`
	Category          string  `json:"category"`
	AvailableQuantity int     `json:"available_quantity"`
}
