# Justfile for fi - File Information CLI Tool

# Default recipe to display help
default:
    @just --list

# Build the binary for current platform
build:
    go build -o fi

# Build with optimizations
build-release:
    go build -ldflags="-s -w" -o fi

# Install locally to /usr/local/bin
install: build
    cp fi /usr/local/bin/fi

# Clean build artifacts
clean:
    rm -f fi
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
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o dist/fi-darwin-amd64
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o dist/fi-darwin-arm64
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o dist/fi-linux-amd64
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o dist/fi-linux-arm64

# Create a release using goreleaser
release:
    goreleaser release --clean --config .github/goreleaser.yml

# Create a snapshot release for testing
snapshot:
    goreleaser release --snapshot --clean --config .github/goreleaser.yml

# Run the binary with example
run FILE="":
    @if [ -z "{{FILE}}" ]; then \
        ./fi --help; \
    else \
        ./fi {{FILE}}; \
    fi

# Show file info with hash
hash FILE:
    ./fi --hash {{FILE}}

# Compare two files
diff FILE1 FILE2:
    ./fi {{FILE1}} {{FILE2}} --diff

# Search for library
lib NAME:
    ./fi --lib {{NAME}}

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
    @if [ -f fi ]; then ls -lh fi | awk '{print $5}'; else echo "Not built yet"; fi

# Development build with race detector
dev:
    go build -race -o fi

# Quick check before commit
check: fmt test lint
    @echo "✓ All checks passed!"
