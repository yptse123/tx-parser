package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"tx-parser/internal/interfaces"
	"tx-parser/internal/storage"
	"tx-parser/pkg/logger"

	"github.com/stretchr/testify/assert"
)

// mockParser is a mock implementation of the Parser interface
type mockParser struct {
	currentBlock int
	subscribed   map[string]bool
	transactions map[string][]interfaces.Transaction
}

func (m *mockParser) GetCurrentBlock() int {
	return m.currentBlock
}

func (m *mockParser) Subscribe(address string) bool {
	if _, exists := m.subscribed[address]; exists {
		return false // Address already subscribed
	}
	m.subscribed[address] = true
	return true // Successfully subscribed
}

func (m *mockParser) GetTransactions(address string) []interfaces.Transaction {
	return m.transactions[address]
}

func TestGetCurrentBlock(t *testing.T) {
	log := logger.GetLogger("debug")
	parser := &mockParser{currentBlock: 123456}
	s := storage.NewMemoryStorage()
	server := NewServer(parser, s, log)

	req, _ := http.NewRequest("GET", "/current-block", nil)
	rr := httptest.NewRecorder()

	server.getCurrentBlock(rr, req)

	// Assert the correct status code and response
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	var response map[string]int
	json.Unmarshal(rr.Body.Bytes(), &response)
	assert.Equal(t, 123456, response["current_block"], "Current block should match the mock response")
}

func TestSubscribe(t *testing.T) {
	log := logger.GetLogger("debug")
	parser := &mockParser{
		subscribed: make(map[string]bool),
	}
	s := storage.NewMemoryStorage()
	server := NewServer(parser, s, log)

	// Test first subscription attempt (success)
	reqBody := []byte(`{"address": "0xNewAddress"}`)
	req, _ := http.NewRequest("POST", "/subscribe", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()

	server.subscribe(rr, req)

	// Assert the correct status code and response
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	var response map[string]string
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Error unmarshalling the response:", err)
	}
	assert.Equal(t, "success", response["status"], "Should subscribe new address successfully")

	// Test second subscription attempt (already subscribed)
	rr = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/subscribe", bytes.NewBuffer(reqBody)) // New request but same address
	req.Header.Set("Content-Type", "application/json")
	server.subscribe(rr, req)

	// Assert the correct status code and response for conflict
	assert.Equal(t, http.StatusConflict, rr.Code, "Status code should be 409")
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal("Error unmarshalling the response:", err)
	}
	assert.Equal(t, "error", response["status"], "Should return an error for already subscribed address")
	assert.Equal(t, "Address already subscribed", response["message"], "Should return appropriate error message for already subscribed address")
}

func TestGetTransactions(t *testing.T) {
	log := logger.GetLogger("debug")
	parser := &mockParser{
		transactions: map[string][]interfaces.Transaction{
			"0xTestAddress": {
				{Hash: "0x1", From: "0xFrom1", To: "0xTestAddress", Value: "100"},
				{Hash: "0x2", From: "0xTestAddress", To: "0xTo1", Value: "200"},
			},
		},
	}
	s := storage.NewMemoryStorage()
	server := NewServer(parser, s, log)

	req, _ := http.NewRequest("GET", "/transactions/0xTestAddress", nil)
	rr := httptest.NewRecorder()

	server.getTransactions(rr, req)

	// Assert the correct status code and response
	assert.Equal(t, http.StatusOK, rr.Code, "Status code should be 200")
	var transactions []interfaces.Transaction
	json.Unmarshal(rr.Body.Bytes(), &transactions)
	assert.Len(t, transactions, 2, "Should return 2 transactions")
	assert.Equal(t, "0x1", transactions[0].Hash, "First transaction hash should match")
	assert.Equal(t, "0xTestAddress", transactions[0].To, "First transaction 'To' address should match")
	assert.Equal(t, "200", transactions[1].Value, "Second transaction value should match")
}

func TestGetTransactions_NoTransactions(t *testing.T) {
	log := logger.GetLogger("debug")
	parser := &mockParser{
		transactions: make(map[string][]interfaces.Transaction),
	}
	s := storage.NewMemoryStorage()
	server := NewServer(parser, s, log)

	req, _ := http.NewRequest("GET", "/transactions/0xTestAddress", nil)
	rr := httptest.NewRecorder()

	server.getTransactions(rr, req)

	// Assert 404 status when no transactions are found
	assert.Equal(t, http.StatusNotFound, rr.Code, "Status code should be 404")
	assert.Equal(t, "No transactions found for the given address\n", rr.Body.String())
}
