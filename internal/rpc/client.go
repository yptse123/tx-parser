package rpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"tx-parser/internal/interfaces"
	"tx-parser/pkg/logger"
)

type Client interface {
	FetchCurrentBlock() (int, error)
	FetchBlockByNumber(int) (*Block, error)
}

type RpcClient struct {
	url string
	log *logger.Logger
}

func NewClient(url string, log *logger.Logger) *RpcClient {
	return &RpcClient{
		url: url,
		log: log,
	}
}

func (c *RpcClient) FetchCurrentBlock() (int, error) {
	payload := `{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`
	resp, err := http.Post(c.url, "application/json", bytes.NewBuffer([]byte(payload)))
	c.log.Info.Printf("request resp: %v", resp)
	if err != nil {
		c.log.Error.Printf("Failed to send request: %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errStr := fmt.Sprintf("failed to fetch current block: status code %d", resp.StatusCode)
		c.log.Error.Printf(errStr)
		return 0, fmt.Errorf(errStr)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.log.Error.Printf("Failed to read response body: %v", err)
		return 0, err
	}

	// Struct to capture the JSON-RPC result
	var result struct {
		Jsonrpc string `json:"jsonrpc"`
		ID      int    `json:"id"`
		Result  string `json:"result"` // The result will be a hexadecimal string
	}

	// Parse the JSON response
	if err := json.Unmarshal(body, &result); err != nil {
		c.log.Error.Printf("Failed to parse JSON response: %v", err)
		return 0, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Trim the "0x" prefix from the hex string if present
	hexStr := strings.TrimPrefix(result.Result, "0x")

	// Convert the hex string to an integer
	blockNumber, err := parseHexToInt(hexStr)
	if err != nil {
		c.log.Error.Printf("Failed to convert block number from hex: %v", err)
		return 0, fmt.Errorf("failed to convert block number from hex: %w", err)
	}

	c.log.Info.Printf("Fetched current block: %d", blockNumber)
	return blockNumber, nil
}

func parseHexToInt(hexStr string) (int, error) {
	var result int
	// Ensure the hex string is prefixed correctly before parsing
	_, err := fmt.Sscanf(hexStr, "%x", &result)
	if err != nil {
		return 0, fmt.Errorf("failed to parse hex string %s: %w", hexStr, err)
	}
	return result, nil
}

// RequestPayload for Ethereum JSON-RPC requests
type RequestPayload struct {
	Jsonrpc string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	Id      int           `json:"id"`
}

// Block and Transaction are used to unmarshal the block data from the RPC
type Block struct {
	Number       string                   `json:"number"`
	Transactions []interfaces.Transaction `json:"transactions"`
}

func (client *RpcClient) FetchBlockByNumber(blockNumber int) (*Block, error) {
	payload := RequestPayload{
		Jsonrpc: "2.0",
		Method:  "eth_getBlockByNumber",
		Params:  []interface{}{fmt.Sprintf("0x%x", blockNumber), true}, // true to include transactions
		Id:      1,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payload: %v", err)
	}

	resp, err := http.Post(client.url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var result struct {
		Jsonrpc string `json:"jsonrpc"`
		Id      int    `json:"id"`
		Result  *Block `json:"result"`
		Error   *struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error,omitempty"`
	}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("RPC error: %s", result.Error.Message)
	}

	return result.Result, nil
}
