APP_NAME=minio-api

.PHONY: run build build-linux tidy swagger

run:
	go run ./cmd/api

build:
	go build -o bin/$(APP_NAME) ./cmd/api

build-linux:
	GOOS=linux GOARCH=amd64 go build -o bin/$(APP_NAME) ./cmd/api

tidy:
	go mod tidy

swagger:
	swag init -g cmd/api/main.go -o docs
