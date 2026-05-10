package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-assistants-catalog/internal/runs/domain"
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

func TestMapErrorAssistantInactive(t *testing.T) {
	apiErr := mapError(domain.ErrAssistantInactive)

	if apiErr.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status: got=%d want=%d", apiErr.StatusCode, http.StatusConflict)
	}
	if apiErr.Code != "ASSISTANT_INACTIVE" {
		t.Fatalf("unexpected code: got=%q want=%q", apiErr.Code, "ASSISTANT_INACTIVE")
	}
}

func TestMapErrorProviderFailed(t *testing.T) {
	apiErr := mapError(domain.ErrProviderFailed)

	if apiErr.StatusCode != http.StatusBadGateway {
		t.Fatalf("unexpected status: got=%d want=%d", apiErr.StatusCode, http.StatusBadGateway)
	}
	if apiErr.Code != "LLM_PROVIDER_ERROR" {
		t.Fatalf("unexpected code: got=%q want=%q", apiErr.Code, "LLM_PROVIDER_ERROR")
	}
}
