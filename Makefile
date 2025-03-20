.PHONY: build test clean run

# Build variables
BINARY_NAME=fetchr
GO=go

build:
	$(GO) build -o $(BINARY_NAME) ./cmd/fetchr

test:
	$(GO) test -v ./...

clean:
	rm -f $(BINARY_NAME)
	$(GO) clean

run:
	$(GO) run ./cmd/fetchr/main.go

fmt:
	$(GO) fmt ./...

vet:
	$(GO) vet ./...

lint:
	golangci-lint run

all: clean build test 