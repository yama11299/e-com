package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/yamadev11/e-com/order/internal/app/spec"
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
