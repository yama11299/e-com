package dl

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yamadev11/e-com/product/internal/app/spec"
)

const (
	productsTable = "products"
)

var (
	selectColumns = []string{
		"id",
		"name",
		"price",
		"category_id",
		"available_quantity",
	}
)

// DL product service data layer interface
type DL interface {
	List(ctx context.Context, request spec.ListRequest) ([]spec.Product, error)
	UpdateQuantity(ctx context.Context, req spec.UpdateQuantityRequest) error
}

type productDL struct {
	db *sqlx.DB
}

// NewProductDL returns product DL client
func NewProductDL(db *sqlx.DB) *productDL {
	return &productDL{db: db}
}

// List returns product list
func (dl *productDL) List(ctx context.Context, request spec.ListRequest) ([]spec.Product, error) {

	response := []spec.Product{}
	q := sq.Select(selectColumns...).From(productsTable)
	if len(request.IDs) > 0 {
		q = q.Where(sq.Eq{"id": request.IDs})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := dl.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		product := spec.Product{}
		err = rows.StructScan(&product)
		if err != nil {
			return nil, err
		}

		response = append(response, product)
	}

	return response, nil
}

// UpdateQuantity updates the quantity for given product id
func (dl *productDL) UpdateQuantity(ctx context.Context, req spec.UpdateQuantityRequest) error {

	q := sq.Update(productsTable).Set("available_quantity", req.Quantity).Where(sq.Eq{"id": req.ID})

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = dl.db.Exec(query, args...)
	if err != nil {
		return err
	}

	log.Printf("quantity updated successfully for the product: %d", req.ID)
	return nil
}
