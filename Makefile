GOCMD=go
BINARY_NAME=flysh

all: clean lint test build

build:
	go build -o bin/$(BINARY_NAME) .

clean:
	rm -rf bin

lint:
	golint ./...

test:
	go test ./...
