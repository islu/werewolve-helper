
build:
	go build -o main

run:
	go run main.go

test:
	go test -v ./...

lint:
	golangci-lint run