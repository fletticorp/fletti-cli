GOCMD=go
BINARY_NAME=fysh

all: clean lint test build

build:
	go build -o bin/$(BINARY_NAME) .

install: build
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	mkdir $(HOME)/.fletaloya
	cp config.yml $(HOME)/.fletaloya

clean:
	rm -rf bin

lint:
	golint ./...

test:
	go test ./...
