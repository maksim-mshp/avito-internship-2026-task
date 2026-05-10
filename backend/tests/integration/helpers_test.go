package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"ai-assistants-catalog/internal/core"
	"ai-assistants-catalog/internal/core/config"

	"github.com/google/uuid"
	testcontainers "github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	postgresImage   = "postgres:18-alpine"
	postgresUser    = "postgres"
	postgresPass    = "postgres"
	jwtToken        = "secret"
	containerWindow = 2 * time.Minute
)

type suite struct {
	baseURL string
	client  *http.Client
}

type response struct {
	statusCode int
	body       []byte
	header     http.Header
}

func newSuite(t *testing.T) *suite {
	t.Helper()

	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}

	ctx, cancel := context.WithTimeout(context.Background(), containerWindow)
	defer cancel()

	databaseName := "integration_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	container, err := tcpostgres.Run(
		ctx,
		postgresImage,
		tcpostgres.BasicWaitStrategies(),
		tcpostgres.WithDatabase(databaseName),
		tcpostgres.WithUsername(postgresUser),
		tcpostgres.WithPassword(postgresPass),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	testcontainers.CleanupContainer(t, container)

	connectionString, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("build postgres connection string: %v", err)
	}

	app := startApp(t, buildConfig(t, connectionString))
	server := httptest.NewServer(app.Server.Handler)

	t.Cleanup(func() {
		server.Close()
		if app.Database != nil {
			app.Database.Close()
		}
	})

	return &suite{
		baseURL: server.URL,
		client:  server.Client(),
	}
}

func buildConfig(t *testing.T, connectionString string) *config.Config {
	t.Helper()

	parsed, err := url.Parse(connectionString)
	if err != nil {
		t.Fatalf("parse connection string %q: %v", connectionString, err)
	}

	password, ok := parsed.User.Password()
	if !ok {
		t.Fatalf("password is missing in connection string %q", connectionString)
	}

	port, err := strconv.Atoi(parsed.Port())
	if err != nil {
		t.Fatalf("parse postgres port from %q: %v", connectionString, err)
	}

	return &config.Config{
		Port: 8080,
		Database: config.Database{
			Host:     parsed.Hostname(),
			Port:     port,
			User:     parsed.User.Username(),
			Password: password,
			Database: strings.TrimPrefix(parsed.Path, "/"),
		},
		JWTToken: jwtToken,
	}
}

func startApp(t *testing.T, cfg *config.Config) *core.App {
	t.Helper()

	app, err := core.Start(cfg)
	if err != nil {
		t.Fatalf("start app: %v", err)
	}

	return app
}

func (s *suite) requestJSON(t *testing.T, method string, path string, body any, token string) response {
	t.Helper()

	var payload io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			t.Fatalf("marshal request body: %v", err)
		}

		payload = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, s.baseURL+path, payload)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		t.Fatalf("perform request %s %s: %v", method, path, err)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			t.Fatalf("close response body: %v", err)
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read response body: %v", err)
	}

	return response{
		statusCode: resp.StatusCode,
		body:       responseBody,
		header:     resp.Header.Clone(),
	}
}

func (s *suite) dummyLogin(t *testing.T, role string) string {
	t.Helper()

	resp := s.requestJSON(t, http.MethodPost, "/dummyLogin", map[string]any{
		"role": role,
	}, "")
	assertStatus(t, resp, http.StatusOK)

	payload := decodeJSON[tokenResponse](t, resp.body)
	if payload.Token == "" {
		t.Fatal("dummy login returned empty token")
	}

	return payload.Token
}

func assertStatus(t *testing.T, resp response, want int) {
	t.Helper()

	if resp.statusCode == want {
		return
	}

	t.Fatalf("unexpected status: got=%d want=%d body=%s", resp.statusCode, want, string(resp.body))
}

func decodeJSON[T any](t *testing.T, data []byte) T {
	t.Helper()

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("decode response body %q: %v", string(data), err)
	}

	return value
}
