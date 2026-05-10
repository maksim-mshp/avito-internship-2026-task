package v1

import (
	"net/http/httptest"
	"testing"
)

func TestParseListQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/runs/my?status=success&page=2&pageSize=30", nil)

	query, apiErr := parseListQuery(req)
	if apiErr != nil {
		t.Fatalf("expected nil error, got %v", apiErr)
	}

	if query.Status == nil || *query.Status != "success" {
		t.Fatalf("unexpected status: %v", query.Status)
	}
	if query.Page != 2 || query.PageSize != 30 {
		t.Fatalf("unexpected pagination: %d %d", query.Page, query.PageSize)
	}
}

func TestParseListQueryDefaults(t *testing.T) {
	req := httptest.NewRequest("GET", "/runs/my", nil)

	query, apiErr := parseListQuery(req)
	if apiErr != nil {
		t.Fatalf("expected nil error, got %v", apiErr)
	}

	if query.Page != 1 || query.PageSize != 20 {
		t.Fatalf("unexpected defaults: %d %d", query.Page, query.PageSize)
	}
}

func TestParseListQueryInvalidPage(t *testing.T) {
	req := httptest.NewRequest("GET", "/runs/my?page=bad", nil)

	_, apiErr := parseListQuery(req)
	if apiErr == nil {
		t.Fatalf("expected error")
	}
}
