SHELL := /bin/bash

.PHONY: check-buf deps generate lint-proto test-core test run db-up db-down

check-buf:
	@command -v buf >/dev/null 2>&1 || (echo "buf is required. macOS: brew install bufbuild/buf/buf" && exit 1)

deps:
	go mod tidy

generate: check-buf
	buf generate

lint-proto: check-buf
	buf lint

# Does not require generated protobuf code.
test-core:
	go test \
		./internal/shared/domain \
		./internal/authentication/domain/... \
		./internal/authentication/application/... \
		./internal/authorization/domain/... \
		./internal/authorization/application/...

# Requires `make generate` first.
test:
	go test ./...

run:
	go run ./cmd/server

db-up:
	docker compose up -d postgres

db-down:
	docker compose down
