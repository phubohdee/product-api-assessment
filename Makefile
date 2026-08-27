include .env
export

.PHONY: run migrate-up migrate-down swagger gen-swagger

run:
	go run ./cmd/api

migrate-up:
	go run ./cmd/api migrate up

migrate-down:
	go run ./cmd/api migrate down

gen-swagger:
	swag init -g cmd/api/main.go -o docs

test:
	go test ./... -v -count=1

test-integration:
	go test ./... -v -count=1 -tags=integration
