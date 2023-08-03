package spec

// ListRequest ...
type ListRequest struct {
	IDs []int `json:"product_ids"`
}

// UpdateQuantityRequest ...
type UpdateQuantityRequest struct {
	ID       int `json:"id"`
	Quantity int `json:"quantity"`
}
