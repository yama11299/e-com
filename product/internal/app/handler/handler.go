package handler

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/yamadev11/e-com/product/internal/app/bl"
)

// List handler to the product List API
func List(svc bl.BL) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// request decoder
		request, err := decodeListRequest(r)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		// call to service
		response, err := svc.List(ctx, request)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		jsonEncodeAPIResponse(ctx, w, response)
	}
}

// UpdateQuantity handler to the product UpdateQuantity API
func UpdateQuantity(svc bl.BL) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()
		// request decoder
		request, err := decodeUpdateQuantityRequest(r)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		// call to service
		err = svc.UpdateQuantity(ctx, request)
		if err != nil {
			jsonEncodeAPIResponse(ctx, w, err)
			return
		}

		jsonEncodeAPIResponse(ctx, w, nil)
	}
}

func jsonEncodeAPIResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	if _, ok := resp.(error); ok {
		w.WriteHeader(http.StatusInternalServerError)
	}

	return json.NewEncoder(w).Encode(resp)
}
