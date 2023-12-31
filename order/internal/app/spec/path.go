package spec

// order service path constant
const (
	CreateOrderPath       = "/orders"
	GetOrderPath          = "/orders/{id:[0-9]+}"
	UpdateOrderStatusPath = "/orders/{id:[0-9]+}"
	CancelOrderPath       = "/orders/{id:[0-9]+}/cancel"
	ReturnOrderPath       = "/orders/{id:[0-9]+}/return"
)
