APP_NAME=minio-api

.PHONY: run build tidy swagger

run:
	go run ./cmd/api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

tidy:
	go mod tidy

swagger:
	swag init -g cmd/api/main.go -o docs
