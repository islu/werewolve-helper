
build:
	go build -o main

run:
	go run main.go

test:
	go test ./... -cover --coverprofile=coverage.out

cover-func: test
	go tool cover -func=coverage.out

cover-html: test
	go tool cover -html=coverage.out

lint:
	golangci-lint run