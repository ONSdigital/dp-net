package request

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const (
	errorDescriptionMalformedRequest = "unable to process request due to a malformed or invalid request body"
)

var (
	errorMalformedRequest = errors.New(errorDescriptionMalformedRequest)
)

// GetJSONRequestBody reads the body from the http.Request and tries to unmarshal it as a JSON to the specified type T
func GetJSONRequestBody[T any](r *http.Request) (*T, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errorMalformedRequest
	}

	var requestBody T
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		return nil, errorMalformedRequest
	}

	return &requestBody, nil
}
