package v1

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"ai-assistants-catalog/internal/assistants/domain"
	"ai-assistants-catalog/internal/core/security"
)

func TestParseListQueryAdminIncludeInactive(t *testing.T) {
	req := httptest.NewRequest("GET", "/assistants?includeInactive=true&page=2&pageSize=20&q=food&tag=recipes", nil)
	req = req.WithContext(security.WithClaims(req.Context(), security.Claims{UserID: "id", Role: security.RoleAdmin}))

	query, apiErr := parseListQuery(req)
	if apiErr != nil {
		t.Fatalf("expected nil error, got %v", apiErr)
	}

	if !query.IncludeInactive {
		t.Fatalf("expected include inactive")
	}
	if query.Page != 2 || query.PageSize != 20 {
		t.Fatalf("unexpected pagination: %d %d", query.Page, query.PageSize)
	}
	if query.Search == nil || *query.Search != "food" {
		t.Fatalf("unexpected search: %v", query.Search)
	}
	if query.Tag == nil || *query.Tag != "recipes" {
		t.Fatalf("unexpected tag: %v", query.Tag)
	}
}

func TestParseListQueryUserIgnoresIncludeInactive(t *testing.T) {
	req := httptest.NewRequest("GET", "/assistants?includeInactive=true", nil)
	req = req.WithContext(security.WithClaims(req.Context(), security.Claims{UserID: "id", Role: security.RoleUser}))

	query, apiErr := parseListQuery(req)
	if apiErr != nil {
		t.Fatalf("expected nil error, got %v", apiErr)
	}

	if query.IncludeInactive {
		t.Fatalf("expected inactive filter to be disabled for user")
	}
}

func TestParseListQueryInvalidBool(t *testing.T) {
	req := httptest.NewRequest("GET", "/assistants?includeInactive=bad", nil)
	req = req.WithContext(security.WithClaims(req.Context(), security.Claims{UserID: "id", Role: security.RoleAdmin}))

	_, apiErr := parseListQuery(req)
	if apiErr == nil {
		t.Fatalf("expected error")
	}
}

func TestMapErrorAssistantNotFound(t *testing.T) {
	apiErr := mapError(domain.ErrNotFound)

	if apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status: got=%d want=%d", apiErr.StatusCode, http.StatusNotFound)
	}
	if apiErr.Code != "ASSISTANT_NOT_FOUND" {
		t.Fatalf("unexpected code: got=%q want=%q", apiErr.Code, "ASSISTANT_NOT_FOUND")
	}
}
