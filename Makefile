# Makefile for Tmux Docker Monitor
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
TIDY=$(GOCMD) mod tidy

# Binary name
BINARY_NAME=tmux-monitor

# Build directory
BUILD_DIR=build

# Main package path
MAIN_PACKAGE_PATH=./cmd/main.go

# Installation directory
INSTALL_DIR=$(HOME)/.tmux/plugins/tmux-monitor

# Build flags
LDFLAGS=-ldflags "-s -w -X main.version=$(shell git describe --tags --always --dirty)"
GCFLAGS=-gcflags="all=-trimpath=$(pwd);-N -l"
ASMFLAGS=-asmflags="all=-trimpath=$(pwd)"

# Optimization flags
OPTIMIZATION_FLAGS=-tags 'osusergo netgo static_build' -installsuffix netgo

# Release flags combining speed, size, and security optimizations
RELEASE_FLAGS=$(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(OPTIMIZATION_FLAGS) -trimpath

.PHONY: all build clean test install release deps

all: deps build

deps:
	$(TIDY)

build: deps
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test: deps
	$(GOTEST) -v ./...

install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	cp tmux-monitor.tmux $(INSTALL_DIR)/
	@echo "Plugin installed to $(INSTALL_DIR)"

release: deps
	CGO_ENABLED=0 $(GOBUILD) $(RELEASE_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)
	@echo "Release version built with optimizations for speed, size, and security"

# Help target
help:
	@echo "Available targets:"
	@echo "  deps     - Get dependencies"
	@echo "  build    - Build the binary"
	@echo "  clean    - Clean build artifacts"
	@echo "  test     - Run tests"
	@echo "  install  - Install the plugin to ~/.tmux/plugins/tmux-docker-monitor"
	@echo "  release  - Build optimized release version (speed, size, and security)"
	@echo "  help     - Show this help message"
