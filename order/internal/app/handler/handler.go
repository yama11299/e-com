package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/yama11299/e-com/order/internal/app/bl"
	"github.com/yama11299/e-com/order/internal/app/bl/dl"
	"github.com/yama11299/e-com/order/internal/app/spec"
)

// CreateOrder handler for CreateOrder API
func CreateOrder(svc bl.BL) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// request decoder
		createOrderRequest, err := decodeCreateOrderRequest(r)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		order, err := svc.Create(ctx, createOrderRequest)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		jsonEncodeAPIResponse(ctx, w, order)
	}
}

// GetOrder handler for GetOrder API
func GetOrder(svc bl.BL) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// request decoder
		orderID, err := decodeGetOrderRequest(r)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		order, err := svc.Get(ctx, orderID)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		jsonEncodeAPIResponse(ctx, w, order)
	}
}

// UpdateOrderStatusHandler handler for UpdateOrderStatus API
func UpdateOrderStatusHandler(svc bl.BL) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// request decoder
		request, err := decodeUpdateOrderStatusRequest(r)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		order, err := svc.UpdateStatus(ctx, request)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		jsonEncodeAPIResponse(ctx, w, order)
	}
}

// CancelOrder handler for CancelOrder API
func CancelOrder(svc bl.BL) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// request decoder
		orderID, err := getOrderID(r)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		request := spec.UpdateOrderStatusRequest{
			OrderID: orderID,
			Status:  dl.Cancelled,
		}
		order, err := svc.CancelAndReturnOrder(ctx, request)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		jsonEncodeAPIResponse(ctx, w, order)
	}
}

// ReturnOrder handler for ReturnOrder API
func ReturnOrder(svc bl.BL) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// request decoder
		orderID, err := getOrderID(r)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		request := spec.UpdateOrderStatusRequest{
			OrderID: orderID,
			Status:  dl.Returned,
		}
		order, err := svc.CancelAndReturnOrder(ctx, request)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		jsonEncodeAPIResponse(ctx, w, order)
	}
}

func jsonEncodeAPIResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if err, ok := resp.(error); ok {
		w.WriteHeader(http.StatusInternalServerError)
		return json.NewEncoder(w).Encode(err.Error())
	}

	return json.NewEncoder(w).Encode(resp)
}
