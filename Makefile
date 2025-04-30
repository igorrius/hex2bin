.PHONY: all build clean install test release

# Variables
BINARY_NAME=hex2bin
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date +%Y%m%d_%H%M%S)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Check for required tools
ZIP_CMD := $(shell command -v zip 2> /dev/null)
TAR_CMD := $(shell command -v tar 2> /dev/null)

# Default target
all: build

# Build for current platform
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf dist/

# Install to GOPATH/bin
install:
	go install $(LDFLAGS)

# Run tests
test:
	go test -v ./...

# Build for all platforms
release: clean
	mkdir -p dist
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64

# Create release archives
release-archives: release
	@echo "Creating release archives..."
	@if [ -n "$(TAR_CMD)" ]; then \
		cd dist && \
		tar -czf $(BINARY_NAME)-linux-amd64-$(VERSION).tar.gz $(BINARY_NAME)-linux-amd64 && \
		tar -czf $(BINARY_NAME)-linux-arm64-$(VERSION).tar.gz $(BINARY_NAME)-linux-arm64 && \
		tar -czf $(BINARY_NAME)-darwin-amd64-$(VERSION).tar.gz $(BINARY_NAME)-darwin-amd64 && \
		tar -czf $(BINARY_NAME)-darwin-arm64-$(VERSION).tar.gz $(BINARY_NAME)-darwin-arm64; \
	else \
		echo "Warning: tar command not found. Skipping tar archives."; \
	fi
	@if [ -n "$(ZIP_CMD)" ]; then \
		cd dist && \
		zip $(BINARY_NAME)-windows-amd64-$(VERSION).zip $(BINARY_NAME)-windows-amd64.exe; \
	else \
		echo "Warning: zip command not found. Skipping Windows zip archive."; \
	fi

# Create checksums
checksums:
	@echo "Creating checksums..."
	@cd dist && \
	if command -v sha256sum >/dev/null 2>&1; then \
		sha256sum * > checksums.txt; \
	elif command -v shasum >/dev/null 2>&1; then \
		shasum -a 256 * > checksums.txt; \
	else \
		echo "Warning: No checksum tool found. Skipping checksums."; \
	fi

# Full release process
full-release: release-archives checksums
	@echo "Release $(VERSION) created in dist/ directory"
	@echo "Files:"
	@ls -l dist/
	@if [ -f dist/checksums.txt ]; then \
		echo "Checksums:"; \
		cat dist/checksums.txt; \
	fi 