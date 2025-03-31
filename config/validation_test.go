package config

import (
	"strings"
	"testing"
)

// TestValidateValidConfig tests validating a valid configuration
func TestValidateValidConfig(t *testing.T) {
	config := Config{
		Server: struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Version     string `yaml:"version"`
		}{
			Name:    "Test Server",
			Version: "1.0.0",
		},
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "API_TOKEN",
		},
		Tools: []ToolConfig{
			{
				Name:        "test-tool",
				Description: "A test tool",
				Endpoint:    "https://example.com/api",
				Method:      "GET",
				Timeout:     30,
				Parameters:  map[string]ParamConfig{},
			},
		},
	}

	err := validateConfig(&config)
	if err != nil {
		t.Errorf("Unexpected error for valid config: %v", err)
	}
}

// TestValidateEmptyToolName tests validating a configuration with an empty tool name
func TestValidateEmptyToolName(t *testing.T) {
	config := Config{
		Tools: []ToolConfig{
			{
				Name:        "", // Empty name
				Description: "A test tool",
				Endpoint:    "https://example.com/api",
				Method:      "GET",
				Timeout:     30,
			},
		},
	}

	err := validateConfig(&config)
	if err == nil {
		t.Error("Expected error for empty tool name, got nil")
		return
	}

	expectedError := "has no name"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestValidateMissingEndpoint tests validating a configuration with a missing endpoint
func TestValidateMissingEndpoint(t *testing.T) {
	config := Config{
		Tools: []ToolConfig{
			{
				Name:        "test-tool",
				Description: "A test tool",
				Endpoint:    "", // Empty endpoint
				Method:      "GET",
				Timeout:     30,
			},
		},
	}

	err := validateConfig(&config)
	if err == nil {
		t.Error("Expected error for missing endpoint, got nil")
		return
	}

	expectedError := "has no endpoint"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestValidateUnsupportedMethod tests validating a configuration with an unsupported HTTP method
func TestValidateUnsupportedMethod(t *testing.T) {
	config := Config{
		Tools: []ToolConfig{
			{
				Name:        "test-tool",
				Description: "A test tool",
				Endpoint:    "https://example.com/api",
				Method:      "DELETE", // Unsupported method
				Timeout:     30,
			},
		},
	}

	err := validateConfig(&config)
	if err == nil {
		t.Error("Expected error for unsupported method, got nil")
		return
	}

	expectedError := "unsupported method"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain '%s', got '%s'", expectedError, err.Error())
	}
}

// TestValidateZeroTimeout tests validating a configuration with a zero timeout value
func TestValidateZeroTimeout(t *testing.T) {
	config := Config{
		Tools: []ToolConfig{
			{
				Name:        "test-tool",
				Description: "A test tool",
				Endpoint:    "https://example.com/api",
				Method:      "GET",
				Timeout:     0, // Zero timeout should be set to default
			},
		},
	}

	err := validateConfig(&config)
	if err != nil {
		t.Errorf("Unexpected error for zero timeout: %v", err)
	}

	if config.Tools[0].Timeout != 30 {
		t.Errorf("Expected default timeout 30, got %d", config.Tools[0].Timeout)
	}
}

// TestValidateNegativeTimeout tests validating a configuration with a negative timeout value
func TestValidateNegativeTimeout(t *testing.T) {
	config := Config{
		Tools: []ToolConfig{
			{
				Name:        "test-tool",
				Description: "A test tool",
				Endpoint:    "https://example.com/api",
				Method:      "GET",
				Timeout:     -10, // Negative timeout should be set to default
			},
		},
	}

	err := validateConfig(&config)
	if err != nil {
		t.Errorf("Unexpected error for negative timeout: %v", err)
	}

	if config.Tools[0].Timeout != 30 {
		t.Errorf("Expected default timeout 30, got %d", config.Tools[0].Timeout)
	}
}
