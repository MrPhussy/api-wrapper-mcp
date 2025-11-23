package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gomcpgo/api_wrapper/config"
	"github.com/gomcpgo/api_wrapper/tool"
	"github.com/mark3labs/mcp-go/mcp"
)

// Session represents an SSE session
type Session struct {
	ID      string
	OutChan chan string // Messages from Server -> Client (SSE)
	Context context.Context
	Cancel  context.CancelFunc
}

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
	ID      interface{}     `json:"id,omitempty"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
	ID      interface{}   `json:"id"`
}

type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func main() {
	// Load configuration
	if len(os.Args) < 2 {
		log.Fatal("Usage: api_wrapper <config.yaml>")
	}

	configFile := os.Args[1]
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create API wrapper handler
	apiToolHandler := tool.NewAPIToolHandler(cfg)

	log.Printf("Starting API Wrapper MCP Server with %d tools...", len(cfg.Tools))

	// Session management
	var (
		mu       sync.Mutex
		sessions = make(map[string]*Session)
	)

	// Healthcheck
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// SSE Endpoint
	http.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
			return
		}

		sessionID := fmt.Sprintf("%d", time.Now().UnixNano())
		ctx, cancel := context.WithCancel(r.Context())

		session := &Session{
			ID:      sessionID,
			OutChan: make(chan string, 10),
			Context: ctx,
			Cancel:  cancel,
		}

		mu.Lock()
		sessions[sessionID] = session
		mu.Unlock()

		defer func() {
			mu.Lock()
			delete(sessions, sessionID)
			mu.Unlock()
			cancel()
			close(session.OutChan)
			log.Printf("Session %s closed", sessionID)
		}()

		// Send endpoint event
		fmt.Fprintf(w, "event: endpoint\ndata: /messages?session_id=%s\n\n", sessionID)
		flusher.Flush()

		log.Printf("New SSE connection: %s", sessionID)

		// Loop to send messages
		for {
			select {
			case msg := <-session.OutChan:
				fmt.Fprintf(w, "event: message\ndata: %s\n\n", msg)
				flusher.Flush()
			case <-ctx.Done():
				return
			}
		}
	})

	// Messages Endpoint
	http.HandleFunc("/messages", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method != "POST" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		sessionID := r.URL.Query().Get("session_id")
		if sessionID == "" {
			http.Error(w, "Missing session_id", http.StatusBadRequest)
			return
		}

		mu.Lock()
		session, exists := sessions[sessionID]
		mu.Unlock()

		if !exists {
			http.Error(w, "Session not found", http.StatusNotFound)
			return
		}

		var req JSONRPCRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		// Handle Request
		go func() {
			resp := handleRequest(req, cfg, apiToolHandler)
			if resp != nil {
				respBytes, _ := json.Marshal(resp)
				select {
				case session.OutChan <- string(respBytes):
				case <-session.Context.Done():
				}
			}
		}()

		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte("Accepted"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting API Wrapper MCP Server on :%s (Custom SSE)", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleRequest(req JSONRPCRequest, cfg *config.Config, handler *tool.APIToolHandler) *JSONRPCResponse {
	switch req.Method {
	case "initialize":
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities": map[string]interface{}{
					"tools": map[string]interface{}{},
				},
				"serverInfo": map[string]string{
					"name":    cfg.Server.Name,
					"version": cfg.Server.Version,
				},
			},
		}
	case "notifications/initialized":
		return nil // No response needed
	case "ping":
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  map[string]string{},
		}
	case "tools/list":
		tools := make([]mcp.Tool, 0, len(cfg.Tools))
		for _, t := range cfg.Tools {
			// Convert config params to JSON schema properties
			props := make(map[string]interface{})
			required := []string{}
			for name, p := range t.Parameters {
				prop := map[string]interface{}{
					"type":        p.Type,
					"description": p.Description,
				}
				if p.Default != nil {
					prop["default"] = p.Default
				}
				if len(p.Enum) > 0 {
					prop["enum"] = p.Enum
				}
				props[name] = prop
				if p.Required {
					required = append(required, name)
				}
			}

			schema := mcp.ToolInputSchema{
				Type:       "object",
				Properties: props,
				Required:   required,
			}

			tools = append(tools, mcp.Tool{
				Name:        t.Name,
				Description: t.Description,
				InputSchema: schema,
			})
		}
		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result: map[string]interface{}{
				"tools": tools,
			},
		}
	case "tools/call":
		var params struct {
			Name      string                 `json:"name"`
			Arguments map[string]interface{} `json:"arguments"`
		}
		if err := json.Unmarshal(req.Params, &params); err != nil {
			return &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &JSONRPCError{
					Code:    -32602,
					Message: "Invalid params",
				},
			}
		}

		// Create mcp.CallToolRequest to pass to handler
		var callReq mcp.CallToolRequest
		callReq.Method = "tools/call"
		callReq.Params.Name = params.Name
		callReq.Params.Arguments = params.Arguments

		result, err := handler.HandleToolCall(context.Background(), callReq)
		if err != nil {
			return &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &JSONRPCError{
					Code:    -32000,
					Message: err.Error(),
				},
			}
		}

		return &JSONRPCResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  result,
		}
	default:
		// Ignore unknown notifications, error on requests
		if req.ID != nil {
			return &JSONRPCResponse{
				JSONRPC: "2.0",
				ID:      req.ID,
				Error: &JSONRPCError{
					Code:    -32601,
					Message: "Method not found",
				},
			}
		}
		return nil
	}
}
