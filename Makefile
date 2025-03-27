.PHONY: build run dev migrate

build:
	go build -o ./cmd/api/main ./cmd/api

run:
	go run ./cmd/api/main.go

# Hot-reloading with CompileDaemon
dev:
<<<<<<< HEAD
	CompileDaemon -build="go build -o ./cmd/api/main ./cmd/api" -command=./cmd/api/main

migrate:
	go run ./cmd/migrate/migrate.go
=======
	CompileDaemon -build="go build -o ./cmd/api/main ./cmd/api" -command=./cmd/api/main
>>>>>>> 96eb9961e0f7697b47c6ea0b2bdd61f4581f4779
