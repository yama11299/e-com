package dl

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/yama11299/e-com/order/internal/app/spec"
)

// Order status constants
const (
	Placed = iota + 1
	Dispatched
	Completed
	Cancelled
	Returned
)

// order table constants
const (
	orderTable      = "orders"
	orderItemsTable = "order_items"
)

// order service constants
var (
	OrderStatusMap = map[int]string{
		Placed:     "Placed",
		Dispatched: "Dispatched",
		Completed:  "Completed",
		Cancelled:  "Cancelled",
		Returned:   "Returned",
	}

	orderTableColumns      = []string{"id", "amount", "discount", "final_amount", "status", "order_date", "dispatch_date"}
	orderItemsTableColumns = []string{"order_id", "product_id", "name", "price", "quantity"}
)

type orderDL struct {
	db *sqlx.DB
}

// NewOrderDL returns order dl client
func NewOrderDL(db *sqlx.DB) *orderDL {
	return &orderDL{db: db}
}

// Create creates the order for provided items
func (dl *orderDL) Create(ctx context.Context, order *spec.Order, items []spec.OrderItem) (*spec.Order, error) {

	currentTime := time.Now()
	y, m, d := time.Now().Date()
	order.OrderDate = fmt.Sprintf("%d/%d/%d", d, m, y)
	order.Status = Placed
	order.FinalAmount = order.Amount - order.Discount

	newOrder := map[string]interface{}{
		"amount":        order.Amount,
		"discount":      order.Discount,
		"final_amount":  order.FinalAmount,
		"status":        order.Status,
		"order_date":    order.OrderDate,
		"dispatch_date": "",
		"created_at":    currentTime.Unix(),
		"updated_at":    currentTime.Unix(),
	}

	q := sq.Insert(orderTable).SetMap(newOrder).Suffix("RETURNING id")
	query, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	tx, err := dl.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	var rowID int
	err = tx.QueryRowx(query, args...).Scan(&rowID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	order.ID = rowID

	// create entries into order item table
	err = dl.CreateOrderItems(ctx, order.ID, items, tx)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return order, nil
}

// CreateOrderItem creates entry for order items
func (dl *orderDL) CreateOrderItems(ctx context.Context, orderID int, items []spec.OrderItem, tx *sqlx.Tx) error {

	q := sq.Insert(orderItemsTable).Columns("order_id", "product_id", "name", "price", "quantity")

	for _, item := range items {
		q = q.Values(orderID, item.ProductID, item.Name, item.Price, item.Quantity)
	}

	query, args, err := q.ToSql()
	if err != nil {
		return err
	}

	_, err = tx.Exec(query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (dl *orderDL) GetOrder(ctx context.Context, orderID int) (spec.Order, error) {

	response := spec.Order{}

	q := sq.Select(orderTableColumns...).From(orderTable).Where(sq.Eq{"id": orderID})
	query, args, err := q.ToSql()
	if err != nil {
		return response, err
	}

	result := dl.db.QueryRowx(query, args...)
	err = result.StructScan(&response)
	if err != nil {
		return response, err
	}

	return response, nil
}

func (dl *orderDL) GetOrderItems(ctx context.Context, orderID int) ([]spec.OrderItem, error) {
	response := []spec.OrderItem{}

	q := sq.Select(orderItemsTableColumns...).From(orderItemsTable).Where(sq.Eq{"order_id": orderID})

	query, args, err := q.ToSql()
	if err != nil {
		return response, err
	}

	rows, err := dl.db.Queryx(query, args...)
	if err != nil {
		return response, err
	}

	orderItem := spec.OrderItem{}
	for rows.Next() {
		err = rows.StructScan(&orderItem)
		if err != nil {
			return response, err
		}

		response = append(response, orderItem)
	}

	return response, nil
}
