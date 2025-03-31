package tool

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gomcpgo/api_wrapper/config"
)

// TestExecuteAPICallGET tests the executeAPICall function with GET requests
func TestExecuteAPICallGET(t *testing.T) {
	// Create a test server to mock HTTP requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/api/test" {
			t.Errorf("Expected path /api/test, got %s", r.URL.Path)
		}

		// Check query parameters
		query := r.URL.Query()
		param1 := query.Get("param1")
		param2 := query.Get("param2")

		if param1 != "value1" {
			t.Errorf("Expected param1=value1, got param1=%s", param1)
		}
		if param2 != "42" {
			t.Errorf("Expected param2=42, got param2=%s", param2)
		}

		// Check authorization header
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Expected Authorization=Bearer test-token, got Authorization=%s", auth)
		}

		// Return a simple response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"data":"test"}`))
	}))
	defer server.Close()

	// Set environment variable for authentication
	os.Setenv("TEST_TOKEN", "test-token")
	defer os.Unsetenv("TEST_TOKEN")

	// Create a test configuration
	cfg := &config.Config{
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
	}

	toolCfg := &config.ToolConfig{
		Name:        "test-get",
		Description: "Test GET endpoint",
		Endpoint:    server.URL + "/api/test",
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
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test executeAPICall with arguments
	ctx := context.Background()
	args := map[string]interface{}{
		"value1": "value1",
		// value2 should be applied from default
	}

	result, err := handler.executeAPICall(ctx, toolCfg, args)
	if err != nil {
		t.Errorf("executeAPICall returned error: %v", err)
	}

	// Check result
	expected := `{"success":true,"data":"test"}`
	if result != expected {
		t.Errorf("Expected result '%s', got '%s'", expected, result)
	}
}

// TestExecuteAPICallPOST tests the executeAPICall function with POST requests
func TestExecuteAPICallPOST(t *testing.T) {
	// Create a test server to mock HTTP requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/api/post" {
			t.Errorf("Expected path /api/post, got %s", r.URL.Path)
		}

		// Check content type
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type=application/json, got Content-Type=%s", contentType)
		}

		// Read body (up to 1024 bytes should be enough for test)
		buf := make([]byte, 1024)
		n, _ := r.Body.Read(buf)
		body := string(buf[:n])

		// Expected JSON body
		expectedBody := `{"key1":"test value","key2":true,"key3":42}`
		if body != expectedBody {
			t.Errorf("Expected body '%s', got '%s'", expectedBody, body)
		}

		// Return a simple response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true,"id":"123"}`))
	}))
	defer server.Close()

	// Create a test configuration
	cfg := &config.Config{
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
	}

	toolCfg := &config.ToolConfig{
		Name:        "test-post",
		Description: "Test POST endpoint",
		Endpoint:    server.URL + "/api/post",
		Method:      "POST",
		Timeout:     30,
		Template:    `{"key1":"{{value1}}","key2":{{value2}},"key3":{{value3}}}`,
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
			"value3": {
				Type:        "number",
				Description: "Third parameter",
				Default:     0,
			},
		},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test executeAPICall with arguments
	ctx := context.Background()
	args := map[string]interface{}{
		"value1": "test value",
		"value2": true,
		"value3": 42,
	}

	result, err := handler.executeAPICall(ctx, toolCfg, args)
	if err != nil {
		t.Errorf("executeAPICall returned error: %v", err)
	}

	// Check result
	expected := `{"success":true,"id":"123"}`
	if result != expected {
		t.Errorf("Expected result '%s', got '%s'", expected, result)
	}
}

// TestExecuteAPICallDefaults tests that default values are applied correctly
func TestExecuteAPICallDefaults(t *testing.T) {
	// Create a test server to mock HTTP requests
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For this test, we're most interested in checking that default values are applied
		query := r.URL.Query()
		param1 := query.Get("param1")
		param2 := query.Get("param2")

		if param1 != "value1" {
			t.Errorf("Expected param1=value1, got param1=%s", param1)
		}
		if param2 != "42" { // Default value
			t.Errorf("Expected param2=42, got param2=%s", param2)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	// Create a test configuration
	cfg := &config.Config{
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
	}

	toolCfg := &config.ToolConfig{
		Name:        "test-defaults",
		Description: "Test defaults",
		Endpoint:    server.URL + "/api/test",
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
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test executeAPICall with only required argument, letting default apply for others
	ctx := context.Background()
	args := map[string]interface{}{
		"value1": "value1",
		// value2 is missing, should use default
	}

	_, err := handler.executeAPICall(ctx, toolCfg, args)
	if err != nil {
		t.Errorf("executeAPICall returned error: %v", err)
	}
}

// TestExecuteAPICallError tests error handling in executeAPICall
func TestExecuteAPICallError(t *testing.T) {
	// Create a test server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"Internal server error"}`))
	}))
	defer server.Close()

	// Create a test configuration
	cfg := &config.Config{
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
	}

	toolCfg := &config.ToolConfig{
		Name:        "test-error",
		Description: "Test error handling",
		Endpoint:    server.URL + "/api/error",
		Method:      "GET",
		Timeout:     30,
		Parameters:  map[string]config.ParamConfig{},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test executeAPICall
	ctx := context.Background()
	args := map[string]interface{}{}

	_, err := handler.executeAPICall(ctx, toolCfg, args)
	if err == nil {
		t.Error("executeAPICall should have returned an error for status 500")
	}
}

// TestExecuteAPICallInvalidURL tests handling of invalid URLs
func TestExecuteAPICallInvalidURL(t *testing.T) {
	// Create a test configuration with invalid URL
	cfg := &config.Config{
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
	}

	toolCfg := &config.ToolConfig{
		Name:        "test-invalid-url",
		Description: "Test invalid URL",
		Endpoint:    "://invalid-url", // Invalid URL
		Method:      "GET",
		Timeout:     30,
		Parameters:  map[string]config.ParamConfig{},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test executeAPICall
	ctx := context.Background()
	args := map[string]interface{}{}

	_, err := handler.executeAPICall(ctx, toolCfg, args)
	if err == nil {
		t.Error("executeAPICall should have returned an error for invalid URL")
	}
}

// TestExecuteAPICallUnsupportedMethod tests handling of unsupported HTTP methods
func TestExecuteAPICallUnsupportedMethod(t *testing.T) {
	// Create a test configuration with unsupported method
	cfg := &config.Config{
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
	}

	toolCfg := &config.ToolConfig{
		Name:        "test-unsupported-method",
		Description: "Test unsupported method",
		Endpoint:    "https://example.com/api/test",
		Method:      "PUT", // Unsupported method
		Timeout:     30,
		Parameters:  map[string]config.ParamConfig{},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test executeAPICall
	ctx := context.Background()
	args := map[string]interface{}{}

	_, err := handler.executeAPICall(ctx, toolCfg, args)
	if err == nil {
		t.Error("executeAPICall should have returned an error for unsupported method")
	}
}

// TestProcessTemplateRegex tests the processTemplateRegex function
func TestProcessTemplateRegex(t *testing.T) {
	// Create a test configuration
	cfg := &config.Config{
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test cases
	testCases := []struct {
		name     string
		template string
		args     map[string]interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "Simple template",
			template: "Hello, {{name}}!",
			args:     map[string]interface{}{"name": "World"},
			expected: "Hello, World!",
			wantErr:  false,
		},
		{
			name:     "JSON template",
			template: `{"name":"{{name}}","value":{{value}}}`,
			args:     map[string]interface{}{"name": "test", "value": 42},
			expected: `{"name":"test","value":42}`,
			wantErr:  false,
		},
		{
			name:     "Template with multiple variables",
			template: `{{var1}}-{{var2}}-{{var3}}`,
			args:     map[string]interface{}{"var1": "a", "var2": "b", "var3": "c"},
			expected: `a-b-c`,
			wantErr:  false,
		},
		{
			name:     "Template with missing variable",
			template: `{{var1}}-{{missing}}-{{var3}}`,
			args:     map[string]interface{}{"var1": "a", "var3": "c"},
			expected: ``,
			wantErr:  true,
		},
		{
			name:     "Template with different data types",
			template: `String: {{string}}, Number: {{number}}, Boolean: {{boolean}}`,
			args: map[string]interface{}{
				"string":  "hello",
				"number":  42,
				"boolean": true,
			},
			expected: `String: hello, Number: 42, Boolean: true`,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handler.processTemplateRegex(tc.template, tc.args)
			
			if tc.wantErr {
				if err == nil {
					t.Errorf("processTemplateRegex() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("processTemplateRegex() error = %v", err)
				}
				if result != tc.expected {
					t.Errorf("processTemplateRegex() = %v, want %v", result, tc.expected)
				}
			}
		})
	}
}

// TestProcessTemplateWithEnvVars tests environment variable handling in templates
func TestProcessTemplateWithEnvVars(t *testing.T) {
	// Set test environment variables
	os.Setenv("TEST_ENV_VAR", "env-value")
	defer os.Unsetenv("TEST_ENV_VAR")

	// Create a test configuration
	cfg := &config.Config{
		Auth: struct {
			TokenEnvVar string `yaml:"token_env_var"`
		}{
			TokenEnvVar: "TEST_TOKEN",
		},
	}

	// Create handler
	handler := NewAPIToolHandler(cfg)

	// Test cases for env vars
	testCases := []struct {
		name     string
		template string
		args     map[string]interface{}
		expected string
		wantErr  bool
	}{
		{
			name:     "Template with env var in args",
			template: `Key: {{api_key}}`,
			args:     map[string]interface{}{"api_key": "{{env:TEST_ENV_VAR}}"},
			expected: `Key: env-value`,
			wantErr:  false,
		},
		{
			name:     "Template with missing env var",
			template: `Key: {{api_key}}`,
			args:     map[string]interface{}{"api_key": "{{env:MISSING_ENV_VAR}}"},
			expected: `Key: `,
			wantErr:  false,
		},
		{
			name:     "Complex template with env var",
			template: `{"auth":"{{token}}","data":"{{value}}"}`,
			args:     map[string]interface{}{
				"token": "{{env:TEST_ENV_VAR}}",
				"value": "test-data",
			},
			expected: `{"auth":"env-value","data":"test-data"}`,
			wantErr:  false,
		},
		{
			name:     "Direct env var in template",
			template: `API Key: {{env:TEST_ENV_VAR}}`,
			args:     map[string]interface{}{},
			expected: `API Key: env-value`,
			wantErr:  false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := handler.processTemplateRegex(tc.template, tc.args)
			
			if tc.wantErr {
				if err == nil {
					t.Errorf("processTemplateRegex() expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("processTemplateRegex() error = %v", err)
				}
				if result != tc.expected {
					t.Errorf("processTemplateRegex() = %v, want %v", result, tc.expected)
				}
			}
		})
	}
}
