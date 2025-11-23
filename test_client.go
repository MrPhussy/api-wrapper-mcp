package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func main() {
	baseURL := "https://notion-mcp-production-a0ef.up.railway.app"
	sseURL := baseURL + "/sse"

	fmt.Printf("Connecting to %s...\n", sseURL)

	// 1. Connect to SSE to get Session ID
	req, _ := http.NewRequest("GET", sseURL, nil)
	req.Header.Set("Accept", "text/event-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to connect to SSE: %v", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	var sessionID string
	var endpoint string

	// Read SSE stream until we get the endpoint event
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("SSE:", line)

		if strings.HasPrefix(line, "data: /messages?session_id=") {
			endpoint = strings.TrimPrefix(line, "data: ")
			parts := strings.Split(endpoint, "=")
			if len(parts) > 1 {
				sessionID = parts[1]
				break
			}
		}
	}

	if sessionID == "" {
		log.Fatal("Failed to get session ID")
	}

	fmt.Printf("Session ID: %s\n", sessionID)
	messagesURL := baseURL + endpoint
	fmt.Printf("Messages URL: %s\n", messagesURL)

	// 2. Send tools/call request
	requestBody := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  "tools/call",
		"params": map[string]interface{}{
			"name": "notion-search",
			"arguments": map[string]interface{}{
				"query": "Speak",
			},
		},
		"id": 1,
	}

	jsonBody, _ := json.Marshal(requestBody)
	postReq, _ := http.NewRequest("POST", messagesURL, bytes.NewBuffer(jsonBody))
	postReq.Header.Set("Content-Type", "application/json")

	fmt.Println("Sending search request...")
	postResp, err := client.Do(postReq)
	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}
	defer postResp.Body.Close()

	if postResp.StatusCode != http.StatusAccepted {
		log.Printf("Unexpected status code: %d", postResp.StatusCode)
	}

	// 3. Read response from SSE stream
	// We need to keep reading the SSE stream from the *original* connection
	// But the scanner loop exited. We need to continue reading.

	fmt.Println("Waiting for response...")
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			// Check if it's the response
			if strings.Contains(data, "result") || strings.Contains(data, "error") {
				fmt.Printf("\nResponse:\n%s\n", data)
				return
			}
		}
	}
}
