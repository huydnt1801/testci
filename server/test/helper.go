package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

type (
	ResultValidator func(t *testing.T)
)

// BuildGetQuery endpoint
func BuildGetQuery(endpoint string, params map[string]string) *http.Request {
	req, _ := http.NewRequest("GET", endpoint, nil)
	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()
	return req
}

// BuildPostQuery end point
func BuildPostQuery(endpoint string, body interface{}) *http.Request {
	b, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	return req
}

// BuildPatchRequest end point
func BuildPatchRequest(endpoint string, body interface{}) *http.Request {
	b, err := json.Marshal(body)
	if err != nil {
		return nil
	}
	req, _ := http.NewRequest("PATCH", endpoint, bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")
	return req
}
