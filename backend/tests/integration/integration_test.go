package integration

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

type tokenResponse struct {
	Token string  `json:"token"`
	User  userDTO `json:"user"`
}

type userDTO struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type registerResponse struct {
	User userDTO `json:"user"`
}

type categoryDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

type assistantDTO struct {
	ID                string   `json:"id"`
	CategoryID        string   `json:"categoryId"`
	CategoryName      *string  `json:"categoryName"`
	Name              string   `json:"name"`
	Description       string   `json:"description"`
	Model             string   `json:"model"`
	SystemPrompt      *string  `json:"systemPrompt"`
	ExampleUserPrompt *string  `json:"exampleUserPrompt"`
	Tags              []string `json:"tags"`
	IsFavorite        bool     `json:"isFavorite"`
	IsActive          bool     `json:"isActive"`
}

type assistantsResponse struct {
	Assistants []assistantDTO `json:"assistants"`
	Pagination paginationDTO  `json:"pagination"`
}

type runDTO struct {
	ID            string  `json:"id"`
	AssistantID   string  `json:"assistantId"`
	AssistantName *string `json:"assistantName"`
	CategoryID    *string `json:"categoryId"`
	CategoryName  *string `json:"categoryName"`
	UserID        string  `json:"userId"`
	Model         string  `json:"model"`
	UserPrompt    string  `json:"userPrompt"`
	Output        *string `json:"output"`
	Status        string  `json:"status"`
	Error         *string `json:"error"`
}

type runsResponse struct {
	Runs       []runDTO      `json:"runs"`
	Pagination paginationDTO `json:"pagination"`
}

type paginationDTO struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Total    int `json:"total"`
}

func TestAdminCreatesCategoryAssistantAndUserRunsAssistant(t *testing.T) {
	s := newSuite(t)

	adminToken := s.dummyLogin(t, "admin")
	userToken := s.dummyLogin(t, "user")

	category := createCategory(t, s, adminToken)
	assistant := createAssistant(t, s, adminToken, category.ID, true)

	listResp := s.requestJSON(t, http.MethodGet, fmt.Sprintf("/assistants?categoryId=%s&page=1&pageSize=10", category.ID), nil, userToken)
	assertStatus(t, listResp, http.StatusOK)

	listBody := decodeJSON[assistantsResponse](t, listResp.body)
	foundAssistant, ok := findAssistant(listBody.Assistants, assistant.ID)
	if !ok {
		t.Fatalf("created assistant %s was not found in user list", assistant.ID)
	}
	if foundAssistant.CategoryID != category.ID {
		t.Fatalf("unexpected assistant category id: got=%q want=%q", foundAssistant.CategoryID, category.ID)
	}
	if foundAssistant.CategoryName == nil || *foundAssistant.CategoryName != category.Name {
		t.Fatalf("unexpected assistant category name: got=%v want=%q", foundAssistant.CategoryName, category.Name)
	}
	if !hasTag(foundAssistant.Tags, "integration") {
		t.Fatalf("expected assistant tags to contain integration, got=%v", foundAssistant.Tags)
	}
	if foundAssistant.IsFavorite {
		t.Fatalf("expected assistant to be not favorite before favorite request")
	}

	tagResp := s.requestJSON(t, http.MethodGet, "/assistants?tag=integration&page=1&pageSize=10", nil, userToken)
	assertStatus(t, tagResp, http.StatusOK)

	tagBody := decodeJSON[assistantsResponse](t, tagResp.body)
	if _, ok = findAssistant(tagBody.Assistants, assistant.ID); !ok {
		t.Fatalf("created assistant %s was not found by tag", assistant.ID)
	}

	favoriteResp := s.requestJSON(t, http.MethodPut, fmt.Sprintf("/assistants/%s/favorite", assistant.ID), nil, userToken)
	assertStatus(t, favoriteResp, http.StatusNoContent)

	favoriteListResp := s.requestJSON(t, http.MethodGet, "/assistants?favoritesOnly=true&page=1&pageSize=10", nil, userToken)
	assertStatus(t, favoriteListResp, http.StatusOK)

	favoriteListBody := decodeJSON[assistantsResponse](t, favoriteListResp.body)
	favoriteAssistant, ok := findAssistant(favoriteListBody.Assistants, assistant.ID)
	if !ok {
		t.Fatalf("created assistant %s was not found in favorites", assistant.ID)
	}
	if !favoriteAssistant.IsFavorite {
		t.Fatalf("expected assistant to be favorite")
	}

	unfavoriteResp := s.requestJSON(t, http.MethodDelete, fmt.Sprintf("/assistants/%s/favorite", assistant.ID), nil, userToken)
	assertStatus(t, unfavoriteResp, http.StatusNoContent)

	emptyFavoriteListResp := s.requestJSON(t, http.MethodGet, "/assistants?favoritesOnly=true&page=1&pageSize=10", nil, userToken)
	assertStatus(t, emptyFavoriteListResp, http.StatusOK)

	emptyFavoriteListBody := decodeJSON[assistantsResponse](t, emptyFavoriteListResp.body)
	if _, ok = findAssistant(emptyFavoriteListBody.Assistants, assistant.ID); ok {
		t.Fatalf("assistant %s was found in favorites after delete", assistant.ID)
	}

	userPrompt := "проверь интеграционный запуск"
	runResp := s.requestJSON(t, http.MethodPost, fmt.Sprintf("/assistants/%s/run", assistant.ID), map[string]any{
		"userPrompt": userPrompt,
	}, userToken)
	assertStatus(t, runResp, http.StatusCreated)

	runBody := decodeJSON[runDTO](t, runResp.body)
	if runBody.ID == "" {
		t.Fatal("run returned empty id")
	}
	if runBody.AssistantID != assistant.ID {
		t.Fatalf("unexpected run assistant id: got=%q want=%q", runBody.AssistantID, assistant.ID)
	}
	if runBody.Status != "success" {
		t.Fatalf("unexpected run status: got=%q want=%q", runBody.Status, "success")
	}
	if runBody.UserPrompt != userPrompt {
		t.Fatalf("unexpected run user prompt: got=%q want=%q", runBody.UserPrompt, userPrompt)
	}
	if runBody.Output == nil || !strings.Contains(*runBody.Output, userPrompt) {
		t.Fatalf("unexpected run output: got=%v", runBody.Output)
	}

	myRunsResp := s.requestJSON(t, http.MethodGet, "/runs/my?page=1&pageSize=10", nil, userToken)
	assertStatus(t, myRunsResp, http.StatusOK)

	myRunsBody := decodeJSON[runsResponse](t, myRunsResp.body)
	myRun, ok := findRun(myRunsBody.Runs, runBody.ID)
	if !ok {
		t.Fatalf("run %s was not found in user history", runBody.ID)
	}
	if myRun.AssistantName == nil || *myRun.AssistantName != assistant.Name {
		t.Fatalf("unexpected run assistant name: got=%v want=%q", myRun.AssistantName, assistant.Name)
	}

	adminRunsResp := s.requestJSON(t, http.MethodGet, fmt.Sprintf("/admin/runs?assistantId=%s&page=1&pageSize=10", assistant.ID), nil, adminToken)
	assertStatus(t, adminRunsResp, http.StatusOK)

	adminRunsBody := decodeJSON[runsResponse](t, adminRunsResp.body)
	if adminRunsBody.Pagination.Total != 1 {
		t.Fatalf("unexpected admin runs total: got=%d want=1", adminRunsBody.Pagination.Total)
	}
	if _, ok = findRun(adminRunsBody.Runs, runBody.ID); !ok {
		t.Fatalf("run %s was not found in admin history", runBody.ID)
	}
}

func TestRegisterAndLogin(t *testing.T) {
	s := newSuite(t)

	email := fmt.Sprintf("integration-%s@example.com", uuid.NewString())
	password := "password"

	registerResp := s.requestJSON(t, http.MethodPost, "/register", map[string]any{
		"email":    email,
		"password": password,
		"role":     "user",
	}, "")
	assertStatus(t, registerResp, http.StatusCreated)

	registerBody := decodeJSON[registerResponse](t, registerResp.body)
	if registerBody.User.ID == "" {
		t.Fatal("register returned empty user id")
	}
	if registerBody.User.Email != email {
		t.Fatalf("unexpected registered email: got=%q want=%q", registerBody.User.Email, email)
	}
	if registerBody.User.Role != "user" {
		t.Fatalf("unexpected registered role: got=%q want=%q", registerBody.User.Role, "user")
	}

	loginResp := s.requestJSON(t, http.MethodPost, "/login", map[string]any{
		"email":    email,
		"password": password,
	}, "")
	assertStatus(t, loginResp, http.StatusOK)

	loginBody := decodeJSON[tokenResponse](t, loginResp.body)
	if loginBody.Token == "" {
		t.Fatal("login returned empty token")
	}
	if loginBody.User.ID != registerBody.User.ID {
		t.Fatalf("unexpected login user id: got=%q want=%q", loginBody.User.ID, registerBody.User.ID)
	}
}

func TestInactiveAssistantIsHiddenAndCannotBeRun(t *testing.T) {
	s := newSuite(t)

	adminToken := s.dummyLogin(t, "admin")
	userToken := s.dummyLogin(t, "user")

	category := createCategory(t, s, adminToken)
	assistant := createAssistant(t, s, adminToken, category.ID, false)

	userListResp := s.requestJSON(t, http.MethodGet, "/assistants?includeInactive=true&page=1&pageSize=10", nil, userToken)
	assertStatus(t, userListResp, http.StatusOK)

	userListBody := decodeJSON[assistantsResponse](t, userListResp.body)
	if _, ok := findAssistant(userListBody.Assistants, assistant.ID); ok {
		t.Fatalf("inactive assistant %s was found in user list", assistant.ID)
	}

	adminListResp := s.requestJSON(t, http.MethodGet, "/assistants?includeInactive=true&page=1&pageSize=10", nil, adminToken)
	assertStatus(t, adminListResp, http.StatusOK)

	adminListBody := decodeJSON[assistantsResponse](t, adminListResp.body)
	if _, ok := findAssistant(adminListBody.Assistants, assistant.ID); !ok {
		t.Fatalf("inactive assistant %s was not found in admin list", assistant.ID)
	}

	runResp := s.requestJSON(t, http.MethodPost, fmt.Sprintf("/assistants/%s/run", assistant.ID), map[string]any{
		"userPrompt": "проверка неактивного ассистента",
	}, userToken)
	assertStatus(t, runResp, http.StatusConflict)
}

func createCategory(t *testing.T, s *suite, token string) categoryDTO {
	t.Helper()

	name := fmt.Sprintf("category-%s", uuid.NewString())
	description := "Integration category"

	resp := s.requestJSON(t, http.MethodPost, "/categories", map[string]any{
		"name":        name,
		"description": description,
	}, token)
	assertStatus(t, resp, http.StatusCreated)

	category := decodeJSON[categoryDTO](t, resp.body)
	if category.ID == "" {
		t.Fatal("create category returned empty id")
	}
	if category.Name != name {
		t.Fatalf("unexpected category name: got=%q want=%q", category.Name, name)
	}

	return category
}

func createAssistant(t *testing.T, s *suite, token string, categoryID string, isActive bool) assistantDTO {
	t.Helper()

	name := fmt.Sprintf("assistant-%s", uuid.NewString())
	description := "Integration assistant"
	model := "mock-smart"
	systemPrompt := "Верни короткий интеграционный ответ"
	exampleUserPrompt := "пример интеграционного контекста"

	resp := s.requestJSON(t, http.MethodPost, "/assistants", map[string]any{
		"categoryId":        categoryID,
		"name":              name,
		"description":       description,
		"model":             model,
		"systemPrompt":      systemPrompt,
		"exampleUserPrompt": exampleUserPrompt,
		"tags":              []string{"integration", "mock"},
		"isActive":          isActive,
	}, token)
	assertStatus(t, resp, http.StatusCreated)

	assistant := decodeJSON[assistantDTO](t, resp.body)
	if assistant.ID == "" {
		t.Fatal("create assistant returned empty id")
	}
	if assistant.CategoryID != categoryID {
		t.Fatalf("unexpected assistant category id: got=%q want=%q", assistant.CategoryID, categoryID)
	}
	if assistant.Name != name {
		t.Fatalf("unexpected assistant name: got=%q want=%q", assistant.Name, name)
	}
	if assistant.IsActive != isActive {
		t.Fatalf("unexpected assistant active flag: got=%v want=%v", assistant.IsActive, isActive)
	}
	if !hasTag(assistant.Tags, "integration") {
		t.Fatalf("unexpected assistant tags: %v", assistant.Tags)
	}

	return assistant
}

func findAssistant(assistants []assistantDTO, assistantID string) (assistantDTO, bool) {
	for _, assistant := range assistants {
		if assistant.ID == assistantID {
			return assistant, true
		}
	}

	return assistantDTO{}, false
}

func hasTag(tags []string, expected string) bool {
	for _, tag := range tags {
		if tag == expected {
			return true
		}
	}

	return false
}

func findRun(runs []runDTO, runID string) (runDTO, bool) {
	for _, run := range runs {
		if run.ID == runID {
			return run, true
		}
	}

	return runDTO{}, false
}
