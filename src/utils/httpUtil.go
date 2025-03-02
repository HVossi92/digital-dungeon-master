package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func SendGetRequest[RespBodyType any](url string) (*RespBodyType, error) {
	// Perform the HTTP POST request
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http GET request failed: %w", err)
	}
	defer response.Body.Close()

	respBody, err := extractResponseBody[RespBodyType](response)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func SendPostRequest[ReqBodyType, RespBodyType any](url string, requestBody ReqBodyType) (*RespBodyType, error) {
	payload, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %w", err)
	}

	response, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("http POST request failed: %w", err)
	}
	defer response.Body.Close()

	respBody, err := extractResponseBody[RespBodyType](response)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func extractResponseBody[RespBodyType any](response *http.Response) (*RespBodyType, error) {
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	// Handle HTTP status codes
	if response.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP error %d: %s", response.StatusCode, string(responseBody))
	}

	// Unmarshal the response body into the specified response type
	var respBody RespBodyType
	if err := json.Unmarshal(responseBody, &respBody); err != nil {
		return nil, fmt.Errorf("error unmarshaling response body: %w, body: %s", err, string(responseBody))
	}

	return &respBody, nil
}
