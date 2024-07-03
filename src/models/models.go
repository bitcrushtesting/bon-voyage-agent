package models

import "encoding/json"

// JSON-RPC request format
type RPCRequest struct {
	Jsonrpc string          `json:"jsonrpc"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
	ID      int             `json:"id"`
}

// JSON-RPC response format
type RPCResponse struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  string `json:"result,omitempty"`
	Error   string `json:"error,omitempty"`
	ID      int    `json:"id"`
}

// RPCError represents a JSON-RPC error
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

type AgentInformation struct {
	Name       string `json:"name"`
	ID         string `json:"id"`
	Version    string `json:"version"`
	ApiVersion string `json:"api_version"`
}

type SetNameMethod struct {
	Name string `json:"name"`
}
