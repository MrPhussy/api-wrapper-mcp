package config

import (
	"strings"
	"testing"
)

// Test validation functions and edge cases
func TestConfigValidation(t *testing.T) {
	// Test simple validation scenarios
	testCases := []struct {
		name          string
		config        Config
		expectError   bool
		errorContains string
	}{
		{
			name: "Valid Config",
			config: Config{
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
			},
			expectError: false,
		},
		{
			name: "Empty Tool Name",
			config: Config{
				Tools: []ToolConfig{
					{
						Name:        "", // Empty name
						Description: "A test tool",
						Endpoint:    "https://example.com/api",
						Method:      "GET",
						Timeout:     30,
					},
				},
			},
			expectError:   true,
			errorContains: "has no name",
		},
		{
			name: "Missing Endpoint",
			config: Config{
				Tools: []ToolConfig{
					{
						Name:        "test-tool",
						Description: "A test tool",
						Endpoint:    "", // Empty endpoint
						Method:      "GET",
						Timeout:     30,
					},
				},
			},
			expectError:   true,
			errorContains: "has no endpoint",
		},
		{
			name: "Unsupported Method",
			config: Config{
				Tools: []ToolConfig{
					{
						Name:        "test-tool",
						Description: "A test tool",
						Endpoint:    "https://example.com/api",
						Method:      "DELETE", // Unsupported method
						Timeout:     30,
					},
				},
			},
			expectError:   true,
			errorContains: "unsupported method",
		},
		{
			name: "Zero Timeout",
			config: Config{
				Tools: []ToolConfig{
					{
						Name:        "test-tool",
						Description: "A test tool",
						Endpoint:    "https://example.com/api",
						Method:      "GET",
						Timeout:     0, // Zero timeout should be set to default
					},
				},
			},
			expectError: false,
		},
		{
			name: "Negative Timeout",
			config: Config{
				Tools: []ToolConfig{
					{
						Name:        "test-tool",
						Description: "A test tool",
						Endpoint:    "https://example.com/api",
						Method:      "GET",
						Timeout:     -10, // Negative timeout should be set to default
					},
				},
			},
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a copy of the config to validate
			cfg := tc.config

			// Validate the configuration
			err := validateConfig(&cfg)

			// Check if error was expected
			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				} else if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("Expected error to contain '%s', got '%s'", tc.errorContains, err.Error())
				}
				return
			}

			// If we're not expecting an error, but got one
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			// Check if defaults were properly applied
			if tc.name == "Zero Timeout" || tc.name == "Negative Timeout" {
				if cfg.Tools[0].Timeout != 30 {
					t.Errorf("Expected default timeout 30, got %d", cfg.Tools[0].Timeout)
				}
			}
		})
	}
}
