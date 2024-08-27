package utils

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/doezaza12/custom-argocd-plugin/pkg/types"
)

func Contains[T comparable](targets []T, comp T) bool {
	for i := 0; i < len(targets); i++ {
		if targets[i] == comp {
			return true
		}
	}
	return false
}

func Request(apiRequest types.ApiRequest) (types.ApiResponse, error) {
	client := &http.Client{}

	allowMethods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}
	if !Contains(allowMethods, apiRequest.Method) {
		return types.ApiResponse{}, errors.New(fmt.Sprintf("method {%s} not allowed.", apiRequest.Method))
	}

	req, err := http.NewRequest(apiRequest.Method, apiRequest.Url, apiRequest.Body)
	if err != nil {
		return types.ApiResponse{}, err
	}

	req.Header = apiRequest.Headers
	res, err := client.Do(req)
	if err != nil {
		return types.ApiResponse{}, err
	}
	defer res.Body.Close()

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		return types.ApiResponse{}, errors.New(fmt.Sprint("Error reading response body:", err))
	}

	apiResponse := types.ApiResponse{}

	apiResponse.Data = resBody
	apiResponse.StatusCode = res.StatusCode

	return apiResponse, nil
}
