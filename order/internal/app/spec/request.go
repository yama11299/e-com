package spec

// CreateOrderRequest ...
type CreateOrderRequest struct {
	Items []OrderItem `json:"order_items"`
}

// UpdateOrderStatusRequest ...
type UpdateOrderStatusRequest struct {
	OrderID int `json:"order_id"`
	Status  int `json:"status"`
}
