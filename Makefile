BINARY_NAME=battery
BIN_DIR=bin
VERSION=1.0.0
LDFLAGS=-ldflags="-s -w -X github.com/twoboots/battery/cmd.Version=$(VERSION)"

.PHONY: all build test test-coverage lint fmt clean install

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build $(LDFLAGS) -o $(BIN_DIR)/$(BINARY_NAME) .

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out

lint:
	go vet ./...

fmt:
	go fmt ./...

install:
	go install $(LDFLAGS) .

clean:
	rm -rf $(BIN_DIR) coverage.out
