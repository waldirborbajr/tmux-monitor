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
LDFLAGS=-ldflags "-s -w"
GCFLAGS=-gcflags="all=-trimpath=$(pwd)"
ASMFLAGS=-asmflags="all=-trimpath=$(pwd)"

# Default build flags for release
RELEASE_FLAGS=$(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) -trimpath

# Optimization flags
SPEED_FLAGS=-tags=speed
SECURITY_FLAGS=-tags=security
SIZE_FLAGS=-tags=size

.PHONY: all build clean test install release-speed release-security release-size deps

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
	cp docker_monitor.tmux $(INSTALL_DIR)/
	cp -r scripts $(INSTALL_DIR)/
	@echo "Plugin installed to $(INSTALL_DIR)"

release-speed: deps
	$(GOBUILD) $(RELEASE_FLAGS) $(SPEED_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

release-security: deps
	$(GOBUILD) $(RELEASE_FLAGS) $(SECURITY_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

release-size: deps
	$(GOBUILD) $(RELEASE_FLAGS) $(SIZE_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

# Target to build all release versions
release-all: release-speed release-security release-size
	@echo "All release versions built"

# Help target
help:
	@echo "Available targets:"
	@echo "  deps             - Get dependencies"
	@echo "  build            - Build the binary"
	@echo "  clean            - Clean build artifacts"
	@echo "  test             - Run tests"
	@echo "  install          - Install the plugin to ~/.tmux/plugins/tmux-docker-monitor"
	@echo "  release-speed    - Build release version optimized for speed"
	@echo "  release-security - Build release version optimized for security"
	@echo "  release-size     - Build release version optimized for small size"
	@echo "  release-all      - Build all release versions"
	@echo "  help             - Show this help message"
