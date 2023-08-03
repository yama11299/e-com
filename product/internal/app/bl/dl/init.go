package dl

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// InitDB initializes and connects to products database
func InitDB() (*sqlx.DB, error) {

	db, err := sqlx.Connect("sqlite3", "products.db")
	if err != nil {
		log.Fatalf("Failed to connect to database, cause: %s", err.Error())
		return nil, err
	}

	createProductsTable := `CREATE TABLE IF NOT EXISTS products (
		id INTEGER PRIMARY KEY ASC,
		name TEXT NOT NULL UNIQUE, 
		price INT NOT NULL,
		category_id INT NOT NULL,
		available_quantity INT NOT NULL,
		created_at INT NOT NULL,
		updated_at INT NOT NULL
	);`

	_, err = db.Exec(createProductsTable)
	if err != nil {
		log.Fatalf("Failed to create products table, cause: %s", err.Error())
		return nil, err
	}

	_ = populateDB(db)

	return db, nil
}

func populateDB(db *sqlx.DB) error {

	populateProductsTable := `INSERT INTO products
	(name, price, available_quantity, category_id, created_at, updated_at) VALUES 
	('Nokia 8.1', 30000, 12, 1, strftime('%s', 'now'), 0),
	('iPhone 14', 65000, 5, 1, strftime('%s', 'now'), 0),
	('Lenovo Ideapad 500', 55000, 7, 2, strftime('%s', 'now'), 0);`

	_, _ = db.Exec(populateProductsTable)

	return nil
}
