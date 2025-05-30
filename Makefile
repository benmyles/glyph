# glyph Makefile

# Binary name
BINARY_NAME=glyph

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOMOD=$(GOCMD) mod
GOGET=$(GOCMD) get

# Build flags
LDFLAGS=-ldflags "-s -w"

.PHONY: all build test clean install help

all: test build

build:
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCMD) clean
	rm -f $(BINARY_NAME)

install:
	@if [ -z "$(DESTDIR)" ]; then \
		echo "Error: DESTDIR not specified. Usage: make install DESTDIR=/path/to/install"; \
		exit 1; \
	fi
	@echo "Installing $(BINARY_NAME) to $(DESTDIR)"
	@mkdir -p $(DESTDIR)
	@rm -f $(DESTDIR)/$(BINARY_NAME)
	@cp -f $(BINARY_NAME) $(DESTDIR)/
	@echo "Installation complete: $(DESTDIR)/$(BINARY_NAME)"

help:
	@echo "Available targets:"
	@echo "  make build    - Build the glyph binary"
	@echo "  make test     - Run all tests"
	@echo "  make install  - Install binary to specified directory"
	@echo "                  Usage: make install DESTDIR=/usr/local/bin"
	@echo "  make clean    - Remove built binary"
	@echo "  make help     - Show this help message"