package dl

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// InitDB initializes and connects to orders database
func InitDB() (*sqlx.DB, error) {

	db, err := sqlx.Connect("sqlite3", "orders.db")
	if err != nil {
		log.Fatalf("Failed to connect to database, cause: %s", err.Error())
		return nil, err
	}

	createOrdersTable := `CREATE TABLE IF NOT EXISTS orders (
		id INTEGER PRIMARY KEY ASC,
		amount INTEGER NOT NULL UNIQUE, 
		discount INTEGER NOT NULL,
		final_amount INTEGER NOT NULL,
		status INTEGER NOT NULL,
		order_date TEXT NOT NULL,
		dispatch_date TEXT NOT NULL,
		created_at INTEGER NOT NULL,
		updated_at INTEGER NOT NULL
	);`

	createOrderItemsTable := `CREATE TABLE IF NOT EXISTS order_items (
		id INTEGER PRIMARY KEY ASC,
		order_id INTEGER NOT NULL,
		product_id INTEGER NOT NULL,
		name TEXT NOT NULL,
		price INTEGER NOT NULL,
		quantity INTEGER NOT NULL
	)`

	_, err = db.Exec(createOrdersTable)
	if err != nil {
		log.Fatalf("Failed to create orders table, cause: %s", err.Error())
		return nil, err
	}

	_, err = db.Exec(createOrderItemsTable)
	if err != nil {
		log.Fatalf("Failed to create order items table, cause: %s", err.Error())
		return nil, err
	}

	return db, nil
}
