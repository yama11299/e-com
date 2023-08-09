package spec

// Order ...
type Order struct {
	ID int `db:"id" json:"id"`
	// Items        []OrderItem `db:"-" json:"items"`
	Amount       float64 `db:"amount" json:"amount"`
	Discount     float64 `db:"discount" json:"discount"`
	FinalAmount  float64 `db:"final_amount" json:"final_amount"`
	Status       int     `db:"status" json:"status"`
	OrderDate    string  `db:"order_date" json:"order_date"`
	DispatchDate string  `db:"dispatch_date" json:"dispatch_date,omitempty"`
	CreatedAt    int     `db:"created_at"`
	UpdatedAt    int     `db:"updated_at"`
}

// OrderItem ...
type OrderItem struct {
	ID        int     `db:"order_id" json:"-"`
	ProductID int     `db:"product_id" json:"product_id"`
	Name      string  `db:"name" json:"name"`
	Price     float64 `db:"price" json:"price"`
	Quantity  int     `db:"quantity" json:"quantity"`
}
