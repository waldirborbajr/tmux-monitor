# VERSION := $(shell cat VERSION | sed 's/^VERSION=//')
# VERSION := $(shell git describe --tags --abbrev=0 | sed 's/^VERSION=//')

# Get the version from git tags, or use a default if not available
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0-dev")

# Remove any "v" prefix and "VERSION=" text if present
VERSION := $(shell echo $(VERSION) | sed -e 's/^v//' -e 's/^VERSION=//')

# If VERSION is empty, set a default development version
ifeq ($(strip $(VERSION)),)
VERSION := 0.0.0-dev
endif

# Makefile for Tmux Docker Monitor
# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean -cache
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
TIDY=$(GOCMD) mod tidy
UPDATE=$(GOCMD) get -u all 

# Binary name
BINARY_NAME=tmux-monitor

# Build directory
BUILD_DIR=build

# Main package path
MAIN_PACKAGE_PATH=./cmd/main.go

# Installation directory
INSTALL_DIR=$(HOME)/.tmux/plugins/tmux-monitor

# Build flags
LDFLAGS=-ldflags "-s -w -X=main.version=$(VERSION)"
GCFLAGS=-gcflags="all=-trimpath=$(pwd);-N -l"
ASMFLAGS=-asmflags="all=-trimpath=$(pwd)"

# Optimization flags
OPTIMIZATION_FLAGS=-tags 'osusergo netgo static_build' -installsuffix netgo

# Release flags combining speed, size, and security optimizations
RELEASE_FLAGS=$(LDFLAGS) $(GCFLAGS) $(ASMFLAGS) $(OPTIMIZATION_FLAGS) -trimpath

.PHONY: all build clean test install release deps

all: deps build

version:
	@echo $(VERSION)

deps:
	$(TIDY)

lint-go:
	golangci-lint cache clean
	golangci-lint run

update:
	$(UPDATE)

build: deps
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE_PATH)

clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

test: deps
	$(GOTEST) -v ./...

install: build
	mkdir -p $(INSTALL_DIR)
	install -D -m 755 $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/tmux-monitor
	install -D -m 755 tmux-monitor.tmux $(INSTALL_DIR)/tmux-monitor.tmux
	# cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	# chmod u+x $(INSTALL_DIR)/tmux-monitor.tmux
	@echo "Plugin installed to $(INSTALL_DIR)"

release: deps
	# Linux AMD64
	env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(RELEASE_FLAGS) -o $(BUILD_DIR)/linux/$(BINARY_NAME) $(MAIN_PACKAGE_PATH) 
	@echo ""
	@echo "Release version built for Linux AMD64 with optimizations for speed, size, and security"
	@echo ""
	# Mac Intel
	env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 $(GOBUILD) $(RELEASE_FLAGS) -o $(BUILD_DIR)/macos/$(BINARY_NAME) $(MAIN_PACKAGE_PATH) 
	@echo ""
	@echo "Release version built for MacOS Intel (AMD64) with optimizations for speed, size, and security"
	@echo ""
	# Mac M1~2
	env GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 $(GOBUILD) $(RELEASE_FLAGS) -o $(BUILD_DIR)/macosarm/$(BINARY_NAME) $(MAIN_PACKAGE_PATH) 
	@echo ""
	@echo "Release version built for MacOS M1~2 (ARM64) with optimizations for speed, size, and security"
	@echo ""
	install -D -m 755 $(BUILD_DIR)/linux/$(BINARY_NAME) $(INSTALL_DIR)/tmux-monitor
	install -D -m 755 tmux-monitor.tmux $(INSTALL_DIR)/tmux-monitor.tmux
	# cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/
	# cp tmux-monitor.tmux $(INSTALL_DIR)/
	# chmod u+x $(INSTALL_DIR)/tmux-monitor.tmux
	@echo ""
	@echo "Installation for Linux AMD64 finished."
	@echo ""

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
