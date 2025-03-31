#!/bin/bash
set -e

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
BIN_DIR="$SCRIPT_DIR/bin"

# Display usage information
function show_usage {
  echo "Usage:"
  echo "  $0 build                  - Build the binary in the bin folder"
  echo "  $0 test [package_path]    - Run tests (optionally for a specific package)"
  echo ""
  echo "Examples:"
  echo "  $0 build                  - Build the api_wrapper binary"
  echo "  $0 test                   - Run all tests"
  echo "  $0 test ./config          - Run tests only for the config package"
}

# Build the binary
function build {
  echo "Building API Wrapper MCP server binary..."
  mkdir -p "$BIN_DIR"
  go build -o "$BIN_DIR/api_wrapper" "$SCRIPT_DIR"
  
  if [ $? -eq 0 ]; then
    echo "Build successful! Binary created at $BIN_DIR/api_wrapper"
  else
    echo "Build failed!"
    exit 1
  fi
}

# Run tests
function run_tests {
  local package_path="$1"
  
  if [ -z "$package_path" ]; then
    echo "Running all tests..."
    go test -v ./...
  else
    echo "Running tests for package $package_path..."
    go test -v "$package_path"
  fi
  
  if [ $? -eq 0 ]; then
    echo "Tests completed successfully!"
  else
    echo "Tests failed!"
    exit 1
  fi
}

# Check command line arguments
if [ $# -lt 1 ]; then
  show_usage
  exit 1
fi

# Process commands
case "$1" in
  build)
    build
    ;;
  test)
    run_tests "$2"
    ;;
  *)
    echo "Error: Unknown command '$1'"
    show_usage
    exit 1
    ;;
esac

exit 0
