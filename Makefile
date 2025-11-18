
build:
	go build -o autonas ./cmd/autonas

run:
	go run ./cmd/autonas run

run-dev:
	ENV=dev go run ./cmd/autonas run

lint:
	golangci-lint run ./...

test:
	go test ./internal/...

test-int:
	(docker rmi autonas:local --force || true) && go test ./integration_tests/... -v -count=1

test-cover:
	go test -coverprofile=coverage.out ./internal/... && \
	go tool cover -html=coverage.out
