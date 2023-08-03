package spec

// Product is the struct for product details
type Product struct {
	ID                int     `db:"id" json:"id"`
	Name              string  `db:"name" json:"name"`
	Price             float64 `db:"price" json:"price"`
	CategoryID        int     `db:"category_id" json:"category_id"`
	AvailableQuantity int     `db:"available_quantity" json:"available_quantity"`
	CreatedAt         uint32  `db:"created_at" json:"created_at"`
	UpdatedAt         uint32  `db:"updated_at" json:"updated_at"`
}
