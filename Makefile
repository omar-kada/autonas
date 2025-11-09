run:
	go run ./cmd/autonas run

run-dev:
	ENV=dev go run ./cmd/autonas run

test:
	go test ./internal/...

test-int:
	(docker rmi autonas:local --force || true) && go test ./integration/... -v -count=1

cover-html:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

build:
	go build -o autonas ./cmd/autonas