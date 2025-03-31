package config

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		name          string
		configFile    string
		expectError   bool
		errorContains string
		validateFunc  func(*testing.T, *Config)
	}{
		{
			name:        "Valid Config",
			configFile:  "testdata/valid_config.yaml",
			expectError: false,
			validateFunc: func(t *testing.T, cfg *Config) {
				// Check server info
				if cfg.Server.Name != "Test API Gateway" {
					t.Errorf("Expected server name 'Test API Gateway', got %s", cfg.Server.Name)
				}
				if cfg.Server.Version != "1.0.0" {
					t.Errorf("Expected server version '1.0.0', got %s", cfg.Server.Version)
				}

				// Check auth
				if cfg.Auth.TokenEnvVar != "TEST_API_TOKEN" {
					t.Errorf("Expected auth token env var 'TEST_API_TOKEN', got %s", cfg.Auth.TokenEnvVar)
				}

				// Check tools count
				if len(cfg.Tools) != 2 {
					t.Errorf("Expected 2 tools, got %d", len(cfg.Tools))
				}

				// Check GET tool
				getToolFound := false
				for _, tool := range cfg.Tools {
					if tool.Name == "test-get" {
						getToolFound = true
						
						if tool.Method != "GET" {
							t.Errorf("Expected GET method, got %s", tool.Method)
						}
						
						if tool.Timeout != 10 {
							t.Errorf("Expected timeout 10, got %d", tool.Timeout)
						}
						
						if len(tool.QueryParams) != 2 {
							t.Errorf("Expected 2 query params, got %d", len(tool.QueryParams))
						}
						
						if len(tool.Parameters) != 2 {
							t.Errorf("Expected 2 parameters, got %d", len(tool.Parameters))
						}
						
						// Check parameter details
						value1Param, ok := tool.Parameters["value1"]
						if !ok {
							t.Error("Expected parameter 'value1' not found")
						} else {
							if !value1Param.Required {
								t.Error("Expected value1 to be required")
							}
							if value1Param.Type != "string" {
								t.Errorf("Expected value1 type to be string, got %s", value1Param.Type)
							}
						}
						
						value2Param, ok := tool.Parameters["value2"]
						if !ok {
							t.Error("Expected parameter 'value2' not found")
						} else {
							if value2Param.Required {
								t.Error("Expected value2 to not be required")
							}
							if value2Param.Type != "number" {
								t.Errorf("Expected value2 type to be number, got %s", value2Param.Type)
							}
							if value2Param.Default != 42.0 {
								t.Errorf("Expected value2 default to be 42, got %v", value2Param.Default)
							}
						}
					}
				}
				if !getToolFound {
					t.Error("Expected tool 'test-get' not found")
				}

				// Check POST tool
				postToolFound := false
				for _, tool := range cfg.Tools {
					if tool.Name == "test-post" {
						postToolFound = true
						
						if tool.Method != "POST" {
							t.Errorf("Expected POST method, got %s", tool.Method)
						}
						
						if tool.Timeout != 20 {
							t.Errorf("Expected timeout 20, got %d", tool.Timeout)
						}
						
						if tool.Template == "" {
							t.Error("Expected non-empty template")
						}
						
						// Check if template contains expected fields
						if !strings.Contains(tool.Template, "{{data}}") {
							t.Error("Expected template to contain '{{data}}'")
						}
						if !strings.Contains(tool.Template, "{{flag}}") {
							t.Error("Expected template to contain '{{flag}}'")
						}
					}
				}
				if !postToolFound {
					t.Error("Expected tool 'test-post' not found")
				}
			},
		},
		{
			name:          "Invalid Method",
			configFile:    "testdata/invalid_method.yaml",
			expectError:   true,
			errorContains: "unsupported method",
		},
		{
			name:          "Missing Endpoint",
			configFile:    "testdata/missing_endpoint.yaml",
			expectError:   true,
			errorContains: "has no endpoint",
		},
		{
			name:        "Default Values",
			configFile:  "testdata/default_values.yaml",
			expectError: false,
			validateFunc: func(t *testing.T, cfg *Config) {
				// Check server defaults
				if cfg.Server.Name != "API Wrapper MCP" {
					t.Errorf("Expected default server name 'API Wrapper MCP', got %s", cfg.Server.Name)
				}
				if cfg.Server.Version != "1.0.0" {
					t.Errorf("Expected default server version '1.0.0', got %s", cfg.Server.Version)
				}

				// Check timeout default
				if len(cfg.Tools) != 1 {
					t.Errorf("Expected 1 tool, got %d", len(cfg.Tools))
				} else {
					tool := cfg.Tools[0]
					if tool.Timeout != 30 {
						t.Errorf("Expected default timeout 30, got %d", tool.Timeout)
					}
				}
			},
		},
		{
			name:          "Non-existent File",
			configFile:    "testdata/nonexistent.yaml",
			expectError:   true,
			errorContains: "failed to read config file",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			configPath := filepath.Join(tc.configFile)
			
			cfg, err := LoadConfig(configPath)
			
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
			
			// Validate configuration
			if tc.validateFunc != nil {
				tc.validateFunc(t, cfg)
			}
		})
	}
}
