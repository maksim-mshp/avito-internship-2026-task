.PHONY: install-deps test-backend test-integration test-frontend lint build openapi coverage

install-deps:
	@cd backend && go install github.com/swaggo/swag/v2/cmd/swag@latest
	@cd backend && go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest
	@cd backend && go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@cd backend && go mod download
	@cd frontend && npm install

test-backend:
	@cd backend && go test -v -short ./...

test-integration:
	@cd backend && go test -v ./tests/integration

test-frontend:
	@cd frontend && npm test

lint:
	@cd backend && files="$$(gofmt -l .)" && test -z "$$files" || (printf '%s\n' "$$files"; exit 1)
	@cd backend && golangci-lint run ./...
	@cd frontend && npm run lint

build:
	@cd backend && go build ./cmd/server
	@cd frontend && npm run build

openapi:
	@cd backend && swag init -g internal/core/http/http.go --output api --outputTypes json,yaml --v3.1
	@cd backend && mv -f api/swagger.json api/openapi.json
	@cd backend && mv -f api/swagger.yaml api/openapi.yml

coverage:
	@cd backend && go test -v -short -coverprofile=coverage.out $$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./... | grep .)
	@cd backend && go tool cover -func=coverage.out
	@cd frontend && npm run test:coverage
