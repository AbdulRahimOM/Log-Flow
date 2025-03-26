.PHONY: build run dev

build:
	go build -o ./cmd/api/main ./cmd/api

run:
	go run ./cmd/api/main.go

# Hot-reloading with CompileDaemon
dev:
	CompileDaemon -build="go build -o ./cmd/api/main ./cmd/api" -command=./cmd/api/main