package rpc

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"tx-parser/pkg/logger"

	"github.com/stretchr/testify/assert"
)

// Mock server for testing RPC responses
func newMockServer(response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, response)
	}))
}

// Test FetchCurrentBlock with a valid response
func TestFetchCurrentBlock(t *testing.T) {
	// Set up a mock server that returns a valid block number in hex format
	mockResponse := `{"jsonrpc":"2.0","id":1,"result":"0xa"}`
	mockServer := newMockServer(mockResponse)
	defer mockServer.Close()

	// Set up the client
	log := logger.GetLogger("debug")
	client := NewClient(mockServer.URL, log)

	// Fetch the current block number
	blockNumber, err := client.FetchCurrentBlock()
	assert.Nil(t, err, "Expected no error when fetching current block")
	assert.Equal(t, 10, blockNumber, "Expected block number to be 10 (0xa in hex)")
}

// Test FetchCurrentBlock with an error response
func TestFetchCurrentBlock_Error(t *testing.T) {
	// Set up a mock server that returns an error response
	mockResponse := `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"Internal error"}}`
	mockServer := newMockServer(mockResponse)
	defer mockServer.Close()

	// Set up the client
	log := logger.GetLogger("debug")
	client := NewClient(mockServer.URL, log)

	// Fetch the current block number (this should return an error)
	blockNumber, err := client.FetchCurrentBlock()
	assert.NotNil(t, err, "Expected an error when fetching current block")
	assert.Equal(t, 0, blockNumber, "Expected block number to be 0 on error")
}

// Test FetchBlockByNumber with a valid block and transactions
func TestFetchBlockByNumber(t *testing.T) {
	// Set up a mock server that returns a block with transactions
	mockResponse := `{
		"jsonrpc": "2.0",
		"id": 1,
		"result": {
			"number": "0x1",
			"transactions": [
				{"hash": "0x1", "from": "0xFrom1", "to": "0xTo1", "value": "100"},
				{"hash": "0x2", "from": "0xFrom2", "to": "0xTo2", "value": "200"}
			]
		}
	}`
	mockServer := newMockServer(mockResponse)
	defer mockServer.Close()

	// Set up the client
	log := logger.GetLogger("debug")
	client := NewClient(mockServer.URL, log)

	// Fetch block by number
	block, err := client.FetchBlockByNumber(1)
	assert.Nil(t, err, "Expected no error when fetching block by number")
	assert.Equal(t, "0x1", block.Number, "Expected block number to be 0x1")
	assert.Len(t, block.Transactions, 2, "Expected 2 transactions in the block")

	// Check the transactions
	assert.Equal(t, "0x1", block.Transactions[0].Hash, "Expected first transaction hash to be 0x1")
	assert.Equal(t, "0x2", block.Transactions[1].Hash, "Expected second transaction hash to be 0x2")
}

// Test FetchBlockByNumber with an error response
func TestFetchBlockByNumber_Error(t *testing.T) {
	// Set up a mock server that returns an error response
	mockResponse := `{"jsonrpc":"2.0","id":1,"error":{"code":-32603,"message":"Internal error"}}`
	mockServer := newMockServer(mockResponse)
	defer mockServer.Close()

	// Set up the client
	log := logger.GetLogger("debug")
	client := NewClient(mockServer.URL, log)

	// Fetch block by number (this should return an error)
	block, err := client.FetchBlockByNumber(1)
	assert.NotNil(t, err, "Expected an error when fetching block by number")
	assert.Nil(t, block, "Expected block to be nil on error")
	assert.True(t, strings.Contains(err.Error(), "Internal error"), "Expected error message to contain 'Internal error'")
}
