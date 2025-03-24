.PHONY: build run dev

build:
	go build -o ./cmd/main ./cmd

run:
	go run ./cmd/main.go

# Hot-reloading with CompileDaemon
dev:
	CompileDaemon -build="go build -o ./cmd/main ./cmd" -command=./cmd/main