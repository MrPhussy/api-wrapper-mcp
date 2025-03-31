package tool

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// processTemplateRegex is a simpler template processor that uses regex for variable substitution
// It supports {{variable}} syntax without requiring dot notation
func (h *APIToolHandler) processTemplateRegex(tmplStr string, args map[string]interface{}) (string, error) {
	// Process env vars in the arguments
	processedArgs := make(map[string]interface{})
	for k, v := range args {
		if strVal, ok := v.(string); ok {
			if strings.HasPrefix(strVal, "{{env:") && strings.HasSuffix(strVal, "}}") {
				envVar := strings.TrimPrefix(strings.TrimSuffix(strVal, "}}"), "{{env:")
				envVal := os.Getenv(envVar)
				processedArgs[k] = envVal
			} else {
				processedArgs[k] = strVal
			}
		} else {
			processedArgs[k] = v
		}
	}

	// Define regex to find {{variable}} patterns
	re := regexp.MustCompile(`{{([^{}]+)}}`)

	// Replace all matches with their values from the args map
	result := re.ReplaceAllStringFunc(tmplStr, func(match string) string {
		// Extract variable name from {{variable}}
		varName := match[2 : len(match)-2]
		
		// Check if this is an environment variable reference
		if strings.HasPrefix(varName, "env:") {
			envVar := strings.TrimPrefix(varName, "env:")
			return os.Getenv(envVar)
		}
		
		// Look up the variable in the args map
		if val, ok := processedArgs[varName]; ok {
			// Convert val to string based on its type
			switch v := val.(type) {
			case string:
				return v
			case int, int64, float64:
				return fmt.Sprintf("%v", v)
			case bool:
				return fmt.Sprintf("%v", v)
			default:
				return fmt.Sprintf("%v", v)
			}
		}
		
		// Variable not found
		return match // Or return an error by using a side-effect
	})

	// Check if any templates were not replaced (which indicates missing variables)
	if re.MatchString(result) {
		// Find which variables were not replaced
		missingVars := re.FindAllString(result, -1)
		return "", fmt.Errorf("missing template variables: %v", missingVars)
	}

	return result, nil
}
