# Justfile for finfo - File Information CLI Tool

# Default recipe to display help
default:
    @just --list

# Build the binary for current platform
build:
    go build -o finfo

# Build with optimizations
build-release:
    go build -ldflags="-s -w" -o finfo

# Install locally to /usr/local/bin
install: build
    cp finfo /usr/local/bin/finfo

# Clean build artifacts
clean:
    rm -f finfo
    rm -rf dist/

# Run tests
test:
    go test ./...

# Run tests with coverage
test-coverage:
    go test -cover ./...

# Run tests with verbose output
test-verbose:
    go test -v ./...

# Format code
fmt:
    go fmt ./...

# Run linter
lint:
    golangci-lint run

# Build for all platforms
build-all:
    mkdir -p dist
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/finfo-darwin-amd64
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/finfo-darwin-arm64
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/finfo-linux-amd64
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/finfo-linux-arm64

# Create a release using goreleaser
release:
    goreleaser release --clean --config .github/goreleaser.yml

# Create a snapshot release for testing
snapshot:
    goreleaser release --snapshot --clean --config .github/goreleaser.yml

# Run the binary with example
run FILE="":
    @if [ -z "{{FILE}}" ]; then \
        ./finfo --help; \
    else \
        ./finfo {{FILE}}; \
    fi

# Show file info with hash
hash FILE:
    ./finfo --hash {{FILE}}

# Compare two files
diff FILE1 FILE2:
    ./finfo {{FILE1}} {{FILE2}} --diff

# Search for library
lib NAME:
    ./finfo --lib {{NAME}}

# Check dependencies
deps:
    go mod download
    go mod verify

# Update dependencies
deps-update:
    go get -u ./...
    go mod tidy

# Show project stats
stats:
    @echo "Lines of code:"
    @wc -l *.go cmd/*.go 2>/dev/null | tail -1
    @echo "\nGo files:"
    @find . -name "*.go" -not -path "./vendor/*" | wc -l
    @echo "\nBinary size (if built):"
    @if [ -f finfo ]; then ls -lh finfo | awk '{print $5}'; else echo "Not built yet"; fi

# Development build with race detector
dev:
    go build -race -o finfo

# Quick check before commit
check: fmt test lint
    @echo "✓ All checks passed!"
