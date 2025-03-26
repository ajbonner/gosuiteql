.PHONY: build test clean install

# Binary name
BINARY_NAME=gosuiteql

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-w -s"

# Build the binary
build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) ./cmd/gosuiteql

# Run tests
test:
	$(GOTEST) -v ./...

# Clean build files
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

# Install dependencies
deps:
	$(GOGET) ./...

# Install the binary to $GOPATH/bin
install:
	$(GOBUILD) $(LDFLAGS) -o $(GOPATH)/bin/$(BINARY_NAME) ./cmd/gosuiteql

# Run the program
run: build
	./$(BINARY_NAME)

# Tidy up go.mod
tidy:
	$(GOMOD) tidy

# Default target
all: deps test build 