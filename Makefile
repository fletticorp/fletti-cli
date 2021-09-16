GOCMD=go
BINARY_NAME=flesh

all: clean lint test build

build:
	go build -o bin/$(BINARY_NAME) .

install: build
	cp bin/$(BINARY_NAME) $(GOPATH)/bin/$(BINARY_NAME)
	mkdir -p $(HOME)/.fletti
	cp -n config.yml $(HOME)/.fletti || true

clean:
	rm -rf bin

lint:
	golint ./...

test:
	go test ./...
