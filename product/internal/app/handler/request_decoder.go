package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/yamadev11/e-com/product/internal/app/spec"
)

func decodeListRequest(r *http.Request) (spec.ListRequest, error) {

	request := spec.ListRequest{}
	productIDs, err := getQueryIDs(r, "product_ids")
	if err != nil {
		return request, err
	}

	request.IDs = productIDs
	return request, nil
}

func decodeUpdateQuantityRequest(r *http.Request) (spec.UpdateQuantityRequest, error) {
	request := spec.UpdateQuantityRequest{}

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&request)
	if err != nil {
		return request, err
	}

	return request, nil
}

func getQueryIDs(r *http.Request, queryParamName string) ([]int, error) {

	var IDs []int
	params := r.URL.Query()
	IDsStr := params.Get(queryParamName)

	IDList := strings.Split(IDsStr, ",")
	for _, param := range IDList {
		param = strings.TrimSpace(param)
		if param != "" {
			id, err := strconv.Atoi(param)
			if err != nil {
				return IDs, fmt.Errorf("invalid query params : %s", queryParamName)
			}
			IDs = append(IDs, id)
		}
	}

	return IDs, nil
}

func getQueryID(r *http.Request, queryParamName string) (int, error) {

	var err error
	var id int
	params := r.URL.Query()
	IDStr := params.Get(queryParamName)
	fmt.Println(IDStr)
	if IDStr != "" {
		id, err = strconv.Atoi(IDStr)
		if err != nil {
			fmt.Println(err)
			return id, fmt.Errorf("invalid query params : %s", queryParamName)
		}
	}

	return id, nil
}
