package bl

import (
	"context"
	"errors"
	"log"

	"github.com/yama11299/e-com/product/internal/app/bl/dl"
	"github.com/yama11299/e-com/product/internal/app/spec"
)

// BL the product service interface
type BL interface {
	List(ctx context.Context, req spec.ListRequest) (spec.ListResponse, error)
	UpdateQuantity(ctx context.Context, req spec.UpdateQuantityRequest) error
}

type bl struct {
	log *log.Logger
	dl  dl.DL
}

// NewProductBL returns the product service client
func NewProductBL(log *log.Logger, dl dl.DL) *bl {
	return &bl{
		log: log,
		dl:  dl,
	}
}

// List returns the product list
func (svc *bl) List(ctx context.Context, req spec.ListRequest) (spec.ListResponse, error) {
	response := spec.ListResponse{}

	productList, err := svc.dl.List(ctx, req)
	if err != nil {
		return response, nil
	}

	response = mapProductModelToListResponse(productList)

	return response, nil
}

func mapProductModelToListResponse(productList []spec.Product) spec.ListResponse {

	listResponse := spec.ListResponse{}

	for _, product := range productList {
		response := spec.GetResponse{
			ID:                product.ID,
			Name:              product.Name,
			Price:             product.Price,
			Category:          dl.CategoryIDNameMap[product.CategoryID],
			AvailableQuantity: product.AvailableQuantity,
		}

		listResponse.Products = append(listResponse.Products, response)
	}

	return listResponse
}

// UpdateQuantity updates the quantity for the given product id
func (svc *bl) UpdateQuantity(ctx context.Context, req spec.UpdateQuantityRequest) error {

	if req.Quantity < 0 {
		return errors.New("invalid quantity")
	}

	return svc.dl.UpdateQuantity(ctx, req)
}
