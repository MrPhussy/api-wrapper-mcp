#!/bin/sh
CONFIG_FILE=${CONFIG_FILE_PATH:-config/example-config.yaml}

echo "Starting MCP Server..."
echo "Config file: $CONFIG_FILE"

# Check if config file exists
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Config file not found at $CONFIG_FILE"
    exit 1
fi

# TODO: Add OpenAPI to YAML conversion logic here if needed in Phase 2

# Run the server
./api_wrapper "$CONFIG_FILE"
