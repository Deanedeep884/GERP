package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"gerp/internal/cli"
	"gerp/internal/mcp"
)

// MCPRequest defines the strict JSON-RPC payload expected from the AI Client.
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

// MCPResponse structures the unified protocol reply back to the Agent.
type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func main() {
	// Ensure the Matrix environments (Spanner/Temporal/BFF) are bound strictly.
	cli.InitConfig()

	// Boot the core MCP STDIO listener.
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Bytes()
		var req MCPRequest
		if err := json.Unmarshal(line, &req); err != nil {
			sendError(nil, -32700, "Parse error", err.Error())
			continue
		}

		handleRequest(req)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "🚨 CRITICAL: MCP STDIO stream ruptured: %v\n", err)
		os.Exit(1)
	}
}

func handleRequest(req MCPRequest) {
	var result interface{}
	var err error

	switch req.Method {
	case "gerp_status":
		result, err = mcp.HandleStatus()
	case "gerp_create_order":
		result, err = mcp.HandleCreateOrder(req.Params)
	case "gerp_audit_view":
		result, err = mcp.HandleAuditView(req.Params)
	case "initialize": // Standard MCP Handshake
		result = map[string]interface{}{
			"protocolVersion": "2024-11-05",
			"serverInfo": map[string]interface{}{
				"name":    "gerp-mcp",
				"version": "1.0.0",
			},
			"capabilities": map[string]interface{}{
				"tools": map[string]interface{}{},
			},
		}
	default:
		sendError(req.ID, -32601, "Method not found", fmt.Sprintf("The GERP Matrix does not recognize the bound method: %s", req.Method))
		return
	}

	if err != nil {
		sendError(req.ID, -32000, "Server error", err.Error())
		return
	}

	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  result,
	}

	out, _ := json.Marshal(resp)
	fmt.Println(string(out))
}

func sendError(id interface{}, code int, message, data string) {
	resp := MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: map[string]interface{}{
			"code":    code,
			"message": message,
			"data":    data,
		},
	}
	out, _ := json.Marshal(resp)
	fmt.Println(string(out))
}
