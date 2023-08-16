package bl

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/yama11299/e-com/order/internal/app/bl/dl"
	"github.com/yama11299/e-com/order/internal/app/spec"
	productGRPC "github.com/yama11299/e-com/product/pb"
)

// DL ...
type DL interface {
	Create(ctx context.Context, order *spec.Order, orderItems []spec.OrderItem) (*spec.Order, error)
	GetOrderItems(ctx context.Context, orderID int) ([]spec.OrderItem, error)
	GetOrder(ctx context.Context, orderID int) (spec.Order, error)
	UpdateStatus(ctx context.Context, req spec.UpdateOrderStatusRequest) error
}

// BL the order service interface
type BL interface {
	Create(ctx context.Context, req spec.CreateOrderRequest) (spec.GetResponse, error)
	Get(ctx context.Context, orderID int) (spec.GetResponse, error)
	UpdateStatus(ctx context.Context, req spec.UpdateOrderStatusRequest) (string, error)
	CancelAndReturnOrder(ctx context.Context, req spec.UpdateOrderStatusRequest) (string, error)
}

type bl struct {
	log     *log.Logger
	dl      DL
	product productGRPC.ProductClient
}

// NewOrderBL returns the order service client
func NewOrderBL(log *log.Logger, dl DL, product productGRPC.ProductClient) *bl {
	return &bl{
		log:     log,
		dl:      dl,
		product: product,
	}
}

// Create creates new order for the provided items
func (svc *bl) Create(ctx context.Context, req spec.CreateOrderRequest) (spec.GetResponse, error) {
	order := &spec.Order{}
	response := spec.GetResponse{}

	ids := make([]int32, len(req.Items))
	for _, item := range req.Items {
		ids = append(ids, int32(item.ProductID))
	}

	listProductResponse, err := svc.product.List(ctx, &productGRPC.ListRequest{Ids: ids})
	if err != nil {
		return response, err
	}

	// create a map for better performance
	productMap := map[int]*productGRPC.GetResponse{}
	for _, product := range listProductResponse.Products {
		productMap[int(product.Id)] = product
	}

	itemMap := map[int]int{}
	for i, item := range req.Items {
		if item.Quantity <= 0 {
			err = errors.New("quantity can not be zero or negative")
			return response, err
		}

		itemMap[item.ProductID] += item.Quantity
		product, ok := productMap[item.ProductID]
		if !ok {
			err = errors.New("invalid ProductID")
			return response, err
		}

		// check whether quantity is within limit or not
		quantity := itemMap[item.ProductID]
		if quantity > int(product.AvailableQuantity) || quantity > 10 {
			err = errors.New("product quantity is beyond limit")
			return response, err
		}

		// filling empty fields
		item.Name = product.Name
		item.Price = float64(product.Price)
		req.Items[i] = item
	}

	var amount, discount float64
	var premiumProductCount int
	for productID, quantity := range itemMap {
		// update product quantity
		product := productMap[productID]
		product.AvailableQuantity -= int32(quantity)

		updateQuantityRequest := &productGRPC.UpdateQuantityRequest{
			Id:       int32(productID),
			Quantity: product.AvailableQuantity,
		}
		_, err = svc.product.UpdateQuantity(ctx, updateQuantityRequest)
		if err != nil {
			return response, err
		}

		amount += float64(quantity) * float64(product.Price)
		if product.Category == "Premium" {
			premiumProductCount++
		}
	}

	// give discount if user has purchased 3 or more premium products
	if premiumProductCount >= 3 {
		discount = amount * 0.1
	}

	order.Amount = amount
	order.Discount = discount
	// create order
	order, err = svc.dl.Create(ctx, order, req.Items)
	if err != nil {
		_ = svc.restoreProductQuantity(ctx, productMap, itemMap)
		return response, err
	}

	response = spec.GetResponse{
		ID:           order.ID,
		Items:        req.Items,
		Amount:       order.Amount,
		Discount:     order.Discount,
		FinalAmount:  order.FinalAmount,
		OrderDate:    order.OrderDate,
		DispatchDate: order.DispatchDate,
		Status:       dl.OrderStatusMap[order.Status],
	}

	return response, nil
}

func (svc *bl) restoreProductQuantity(ctx context.Context, products map[int]*productGRPC.GetResponse, items map[int]int) error {

	for productID, quantity := range items {
		// update product quantity
		product := products[productID]
		product.AvailableQuantity += int32(quantity)

		updateQuantityRequest := &productGRPC.UpdateQuantityRequest{
			Id:       int32(productID),
			Quantity: product.AvailableQuantity,
		}
		_, err := svc.product.UpdateQuantity(ctx, updateQuantityRequest)
		if err != nil {
			return err
		}
	}

	return nil
}

func (svc *bl) Get(ctx context.Context, orderID int) (spec.GetResponse, error) {

	var response spec.GetResponse
	order, err := svc.dl.GetOrder(ctx, orderID)
	if err != nil {
		return response, err
	}

	orderItems, err := svc.dl.GetOrderItems(ctx, orderID)
	if err != nil {
		return response, err
	}

	response = spec.GetResponse{
		ID:           order.ID,
		Items:        orderItems,
		Amount:       order.Amount,
		Discount:     order.Discount,
		FinalAmount:  order.FinalAmount,
		Status:       dl.OrderStatusMap[order.Status],
		OrderDate:    order.OrderDate,
		DispatchDate: order.DispatchDate,
	}

	return response, nil
}

// UpdateStatus updates order status
func (svc *bl) UpdateStatus(ctx context.Context, req spec.UpdateOrderStatusRequest) (string, error) {
	var response string

	// get order details
	order, err := svc.dl.GetOrder(ctx, req.OrderID)
	if err != nil {
		return response, err
	}

	// validate status
	if order.Status == dl.Cancelled || order.Status == dl.Returned || order.Status >= req.Status {
		return response, errors.New("couldn't update order status, invalid status provided")
	}

	// update status
	err = svc.dl.UpdateStatus(ctx, req)
	if err != nil {
		return response, err
	}

	response = fmt.Sprintf("updated status as %s for order id %d", dl.OrderStatusMap[req.Status], req.OrderID)
	return response, nil
}

// CancelAndReturnOrder sets order status as cancelled or returned
func (svc *bl) CancelAndReturnOrder(ctx context.Context, req spec.UpdateOrderStatusRequest) (string, error) {
	var response string

	// get order details
	order, err := svc.dl.GetOrder(ctx, req.OrderID)
	if err != nil {
		return response, err
	}

	// validate status
	if req.Status == dl.Returned && order.Status != dl.Delivered {
		return response, errors.New("couldn't update order status, invalid status provided")
	}

	if order.Status == dl.Cancelled || order.Status >= req.Status {
		return response, errors.New("couldn't update order status, invalid status provided")
	}

	// get order items to update their available quantity in product service
	orderItems, err := svc.dl.GetOrderItems(ctx, order.ID)
	if err != nil {
		return response, err
	}

	ids := make([]int32, len(orderItems))
	for _, item := range orderItems {
		ids = append(ids, int32(item.ProductID))
	}

	listProductResponse, err := svc.product.List(ctx, &productGRPC.ListRequest{Ids: ids})
	if err != nil {
		return response, err
	}

	// create a map for better performance
	productMap := map[int]*productGRPC.GetResponse{}
	for _, product := range listProductResponse.Products {
		productMap[int(product.Id)] = product
	}

	for _, item := range orderItems {
		product := productMap[item.ProductID]
		product.AvailableQuantity += int32(item.Quantity)

		updateQuantityReq := &productGRPC.UpdateQuantityRequest{
			Id:       product.Id,
			Quantity: product.AvailableQuantity,
		}

		_, err = svc.product.UpdateQuantity(ctx, updateQuantityReq)
		if err != nil {
			return response, err
		}
	}

	// update status
	updateStatusReq := spec.UpdateOrderStatusRequest{
		OrderID: order.ID,
		Status:  req.Status,
	}

	err = svc.dl.UpdateStatus(ctx, updateStatusReq)
	if err != nil {
		return response, err
	}

	response = fmt.Sprintf("updated status as %s for order id %d",
		dl.OrderStatusMap[updateStatusReq.Status], updateStatusReq.OrderID)

	return response, nil
}
