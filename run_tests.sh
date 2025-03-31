#!/bin/bash
set -e

# Run tests
echo "Running API Wrapper MCP server tests..."

# Run all tests with verbose output
go test -v ./...

echo "Tests completed successfully!"
