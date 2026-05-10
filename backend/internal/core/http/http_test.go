package http

import (
	"net/http/httptest"
	"strings"
	"testing"
)

type parseRequest struct {
	Name string `json:"name"`
}

func TestParseJSONBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "application/json")

	var body parseRequest
	apiErr := ParseJSONBody(req, &body)
	if apiErr != nil {
		t.Fatalf("expected nil error, got %v", apiErr)
	}
	if body.Name != "test" {
		t.Fatalf("unexpected name: %s", body.Name)
	}
}

func TestParseJSONBodyInvalidContentType(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "text/plain")

	var body parseRequest
	apiErr := ParseJSONBody(req, &body)
	if apiErr == nil || apiErr.Code != ErrInvalidRequest.Code {
		t.Fatalf("expected invalid request, got %v", apiErr)
	}
}

func TestParseJSONBodyInvalidJSON(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{`))
	req.Header.Set("Content-Type", "application/json")

	var body parseRequest
	apiErr := ParseJSONBody(req, &body)
	if apiErr == nil || apiErr.Code != ErrInvalidRequest.Code {
		t.Fatalf("expected invalid request, got %v", apiErr)
	}
}

func TestNewServerInvalidPort(t *testing.T) {
	_, err := NewServer(0, nil)
	if err == nil {
		t.Fatalf("expected error")
	}
}
