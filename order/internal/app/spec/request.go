package spec

// CreateOrderRequest ...
type CreateOrderRequest struct {
	Items []OrderItem `json:"order_items"`
}
