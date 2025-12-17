
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
	(docker stop autonas || true) && \
	(docker rm autonas || true) && \
	(docker rmi autonas:local --force || true) && \
	cd backend && go test ./integration_tests/... -v -count=1

test-cover:
	cd backend && go test -coverprofile=coverage.out ./internal/... && \
	go tool cover -html=coverage.out

gen-api : 
	make tsp-gen 
	make oapi-gen
	make orval-gen

tsp-gen:
	cd api && npx tsp compile .

oapi-gen:
	cd backend && go tool oapi-codegen -config oapi-codegen.yaml ../api/tsp-output/schema/openapi.1.0.yaml

orval-gen:
	cd frontend && npx orval --config orval.config.js

# docker-run:
# 	docker compose --env-file ./_ignore_.env up --build