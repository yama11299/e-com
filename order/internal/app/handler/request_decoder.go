package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/yama11299/e-com/order/internal/app/spec"
)

func decodeCreateOrderRequest(r *http.Request) (spec.CreateOrderRequest, error) {

	request := spec.CreateOrderRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		fmt.Println("failed to decode request", err.Error())
		return request, err
	}

	return request, nil
}

func decodeGetOrderRequest(r *http.Request) (int, error) {

	orderID, err := getOrderID(r)
	if err != nil {
		return 0, err
	}

	return orderID, nil
}

func getOrderID(r *http.Request) (int, error) {

	vars := mux.Vars(r)
	orderIDStr, exists := vars["id"]
	if !exists {
		err := errors.New("key not found")
		err = errors.Join(err, errors.New("msg=failed to get value for key order id"))
		return -1, err
	}

	orderID, conversionError := strconv.ParseInt(orderIDStr, 10, 32)
	if conversionError != nil || orderID < 0 {
		err := errors.New("invalid order id")
		return -1, err
	}

	return int(orderID), nil
}
