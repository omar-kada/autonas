run:
	go run ./cmd/autonas run

run-dev:
	ENV=dev go run ./cmd/autonas run

test:
	go test ./...

cover-html:
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

build:
	go build -o autonas ./cmd/autonas