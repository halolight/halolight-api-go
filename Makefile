.PHONY: dev build run test clean

dev:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

run: build
	./bin/server

test:
	go test ./...

clean:
	rm -rf bin/
