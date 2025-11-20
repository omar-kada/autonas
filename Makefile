
build:
	cd backend && go build -o autonas ./cmd/autonas

run:
	cd backend && go run ./cmd/autonas run

run-dev:
	ENV=dev cd backend && go run ./cmd/autonas run

lint:
	cd backend && golangci-lint run ./...

test:
	cd backend && go test ./internal/...

test-int:
	(docker rmi autonas:local --force || true) && cd backend && go test ./integration_tests/... -v -count=1

test-cover:
	cd backend && go test -coverprofile=coverage.out ./internal/... && \
	go tool cover -html=coverage.out
