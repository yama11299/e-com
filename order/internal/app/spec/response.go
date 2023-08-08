package spec

// GetResponse ...
type GetResponse struct {
	ID           int         `json:"order_id"`
	Items        []OrderItem `json:"order_items"`
	Amount       float64     `json:"amount"`
	Discount     float64     `json:"discount"`
	FinalAmount  float64     `json:"finalAmount"`
	Status       string      `json:"status"`
	OrderDate    string      `json:"orderDate"`
	DispatchDate string      `json:"dispatchDate,omitempty"`
}
