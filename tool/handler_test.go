package tool

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/gomcpgo/api_wrapper/config"
)

// TestListTools tests the ListTools function
func TestListTools(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		Server: struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Version     string `yaml:"version"`
		}{
			Name:        "Test Server",
			Description: "Test server for unit tests",
			Version:     "1.0.0",
		},
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
		Tools: []config.ToolConfig{
			{
				Name:        "test-get",
				Description: "Test GET tool",
				Endpoint:    "https://example.com/api/get",
				Method:      "GET",
				Timeout:     30,
				QueryParams: map[string]string{
					"param1": "{{value1}}",
					"param2": "{{value2}}",
				},
				Parameters: map[string]config.ParamConfig{
					"value1": {
						Type:        "string",
						Description: "First parameter",
						Required:    true,
					},
					"value2": {
						Type:        "number",
						Description: "Second parameter",
						Default:     42,
					},
				},
			},
			{
				Name:        "test-post",
				Description: "Test POST tool",
				Endpoint:    "https://example.com/api/post",
				Method:      "POST",
				Timeout:     60,
				Template:    `{"key1":"{{value1}}","key2":{{value2}}}`,
				Parameters: map[string]config.ParamConfig{
					"value1": {
						Type:        "string",
						Description: "First parameter",
						Required:    true,
					},
					"value2": {
						Type:        "boolean",
						Description: "Second parameter",
						Default:     false,
					},
				},
			},
		},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test ListTools
	ctx := context.Background()
	resp, err := handler.ListTools(ctx)

	// Check results
	if err != nil {
		t.Errorf("ListTools returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("ListTools returned nil response")
	}

	// Check number of tools
	if len(resp.Tools) != 2 {
		t.Errorf("Expected 2 tools, got %d", len(resp.Tools))
	}

	// Check the properties of each tool
	for _, tool := range resp.Tools {
		switch tool.Name {
		case "test-get":
			if tool.Description != "Test GET tool" {
				t.Errorf("Expected description 'Test GET tool', got '%s'", tool.Description)
			}
			// Check the schema
			verifyToolSchema(t, tool.InputSchema, []string{"value1", "value2"}, []string{"value1"})

		case "test-post":
			if tool.Description != "Test POST tool" {
				t.Errorf("Expected description 'Test POST tool', got '%s'", tool.Description)
			}
			// Check the schema
			verifyToolSchema(t, tool.InputSchema, []string{"value1", "value2"}, []string{"value1"})

		default:
			t.Errorf("Unexpected tool name: %s", tool.Name)
		}
	}
}

// verifyToolSchema checks the JSON schema of a tool
func verifyToolSchema(t *testing.T, schemaJSON json.RawMessage, expectedParams []string, expectedRequired []string) {
	t.Helper()

	// Parse the schema
	var schema struct {
		Type       string                 `json:"type"`
		Properties map[string]interface{} `json:"properties"`
		Required   []string               `json:"required"`
	}

	if err := json.Unmarshal(schemaJSON, &schema); err != nil {
		t.Errorf("Failed to parse schema: %v", err)
		return
	}

	// Check schema type
	if schema.Type != "object" {
		t.Errorf("Expected schema type 'object', got '%s'", schema.Type)
	}

	// Check properties
	if len(schema.Properties) != len(expectedParams) {
		t.Errorf("Expected %d properties, got %d", len(expectedParams), len(schema.Properties))
	}

	for _, paramName := range expectedParams {
		if _, ok := schema.Properties[paramName]; !ok {
			t.Errorf("Expected property '%s' not found in schema", paramName)
		}
	}

	// Check required fields
	if len(schema.Required) != len(expectedRequired) {
		t.Errorf("Expected %d required fields, got %d", len(expectedRequired), len(schema.Required))
	}

	for _, req := range expectedRequired {
		found := false
		for _, actual := range schema.Required {
			if actual == req {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected required field '%s' not found in required list", req)
		}
	}
}

// TestListToolsEmpty tests the ListTools function with no tools
func TestListToolsEmpty(t *testing.T) {
	// Create a test configuration with no tools
	cfg := &config.Config{
		Server: struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Version     string `yaml:"version"`
		}{
			Name:        "Empty Server",
			Description: "Server with no tools",
			Version:     "1.0.0",
		},
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
		Tools: []config.ToolConfig{},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test ListTools
	ctx := context.Background()
	resp, err := handler.ListTools(ctx)

	// Check results
	if err != nil {
		t.Errorf("ListTools returned error: %v", err)
	}

	if resp == nil {
		t.Fatal("ListTools returned nil response")
	}

	// Check number of tools
	if len(resp.Tools) != 0 {
		t.Errorf("Expected 0 tools, got %d", len(resp.Tools))
	}
}

// TestToolSchema tests schema generation for different parameter types
func TestToolSchema(t *testing.T) {
	// Create a test configuration with different parameter types
	cfg := &config.Config{
		Server: struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Version     string `yaml:"version"`
		}{
			Name:        "Test Server",
			Description: "Test server for schema validation",
			Version:     "1.0.0",
		},
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
		Tools: []config.ToolConfig{
			{
				Name:        "test-types",
				Description: "Tool with different parameter types",
				Endpoint:    "https://example.com/api/test",
				Method:      "POST",
				Timeout:     30,
				Template:    `{}`,
				Parameters: map[string]config.ParamConfig{
					"string_param": {
						Type:        "string",
						Description: "String parameter",
						Required:    true,
					},
					"number_param": {
						Type:        "number",
						Description: "Number parameter",
						Default:     123,
					},
					"boolean_param": {
						Type:        "boolean",
						Description: "Boolean parameter",
						Default:     true,
					},
					"enum_param": {
						Type:        "string",
						Description: "Enum parameter",
						Enum:        []string{"option1", "option2", "option3"},
					},
				},
			},
		},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test ListTools
	ctx := context.Background()
	resp, err := handler.ListTools(ctx)

	// Check results
	if err != nil {
		t.Errorf("ListTools returned error: %v", err)
	}

	if len(resp.Tools) != 1 {
		t.Fatalf("Expected 1 tool, got %d", len(resp.Tools))
	}

	tool := resp.Tools[0]
	
	// Parse the schema
	var schema struct {
		Type       string                       `json:"type"`
		Properties map[string]json.RawMessage   `json:"properties"`
		Required   []string                     `json:"required"`
	}

	if err := json.Unmarshal(tool.InputSchema, &schema); err != nil {
		t.Fatalf("Failed to parse schema: %v", err)
	}

	// Check string parameter
	var stringParam struct {
		Type        string `json:"type"`
		Description string `json:"description"`
	}
	if err := json.Unmarshal(schema.Properties["string_param"], &stringParam); err != nil {
		t.Errorf("Failed to parse string parameter: %v", err)
	}
	if stringParam.Type != "string" {
		t.Errorf("Expected string_param type to be 'string', got '%s'", stringParam.Type)
	}
	if stringParam.Description != "String parameter" {
		t.Errorf("Expected string_param description to be 'String parameter', got '%s'", stringParam.Description)
	}

	// Check number parameter
	var numberParam struct {
		Type        string      `json:"type"`
		Description string      `json:"description"`
		Default     json.Number `json:"default"`
	}
	if err := json.Unmarshal(schema.Properties["number_param"], &numberParam); err != nil {
		t.Errorf("Failed to parse number parameter: %v", err)
	}
	if numberParam.Type != "number" {
		t.Errorf("Expected number_param type to be 'number', got '%s'", numberParam.Type)
	}

	// Check boolean parameter
	var boolParam struct {
		Type        string `json:"type"`
		Description string `json:"description"`
		Default     bool   `json:"default"`
	}
	if err := json.Unmarshal(schema.Properties["boolean_param"], &boolParam); err != nil {
		t.Errorf("Failed to parse boolean parameter: %v", err)
	}
	if boolParam.Type != "boolean" {
		t.Errorf("Expected boolean_param type to be 'boolean', got '%s'", boolParam.Type)
	}
	if !boolParam.Default {
		t.Errorf("Expected boolean_param default to be true, got %v", boolParam.Default)
	}

	// Check enum parameter
	var enumParam struct {
		Type        string   `json:"type"`
		Description string   `json:"description"`
		Enum        []string `json:"enum"`
	}
	if err := json.Unmarshal(schema.Properties["enum_param"], &enumParam); err != nil {
		t.Errorf("Failed to parse enum parameter: %v", err)
	}
	if enumParam.Type != "string" {
		t.Errorf("Expected enum_param type to be 'string', got '%s'", enumParam.Type)
	}
	if len(enumParam.Enum) != 3 {
		t.Errorf("Expected enum_param to have 3 options, got %d", len(enumParam.Enum))
	}

	// Check required fields
	foundRequired := false
	for _, req := range schema.Required {
		if req == "string_param" {
			foundRequired = true
			break
		}
	}
	if !foundRequired {
		t.Errorf("Expected 'string_param' to be in required fields: %v", schema.Required)
	}
}
