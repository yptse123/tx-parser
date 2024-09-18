package parser

import (
	"testing"
	"tx-parser/internal/interfaces"
	"tx-parser/internal/rpc"
	"tx-parser/internal/storage"
	"tx-parser/pkg/logger"

	"github.com/stretchr/testify/assert"
)

// Mock implementation of the rpc.Client
type mockRPCClient struct{}

func (m *mockRPCClient) FetchCurrentBlock() (int, error) {
	return 10, nil
}

func (m *mockRPCClient) FetchBlockByNumber(blockNumber int) (*rpc.Block, error) {
	// Return mock block with transactions based on the block number
	switch blockNumber {
	case 1:
		return &rpc.Block{
			Transactions: []interfaces.Transaction{
				{Hash: "0x1", From: "0xfrom1", To: "0xtestaddress", Value: "100"},
				{Hash: "0x2", From: "0xtestaddress", To: "0xto1", Value: "200"},
			},
		}, nil
	case 2:
		return &rpc.Block{
			Transactions: []interfaces.Transaction{
				{Hash: "0x3", From: "0xtestaddress", To: "0xto2", Value: "300"},
			},
		}, nil
	default:
		return &rpc.Block{
			Transactions: []interfaces.Transaction{
				{Hash: "0x4", From: "0xfrom2", To: "0xanotheraddress", Value: "400"},
			},
		}, nil
	}
}

// Test fetching current block during initialization
func TestNewEthParser(t *testing.T) {
	log := logger.GetLogger("debug")
	client := &mockRPCClient{}
	storage := storage.NewMemoryStorage()

	parser := NewEthParser(client, storage, log)
	assert.Equal(t, 10, parser.currentBlock, "The initial block number should be 10")
}

// Test subscribing to an address
func TestSubscribe(t *testing.T) {
	log := logger.GetLogger("debug")
	client := &mockRPCClient{}
	storage := storage.NewMemoryStorage()

	parser := NewEthParser(client, storage, log)

	// Test subscribing to a new address
	address := "0xTestAddress"
	subscribed := parser.Subscribe(address)
	assert.True(t, subscribed, "Should successfully subscribe a new address")

	// Test subscribing again to the same address
	subscribed = parser.Subscribe(address)
	assert.False(t, subscribed, "Should not subscribe the same address again")
}

// Test fetching transactions for an address
func TestGetTransactions(t *testing.T) {
	log := logger.GetLogger("debug")
	client := &mockRPCClient{}
	mockStorage := storage.NewMemoryStorage()

	parser := NewEthParser(client, mockStorage, log)

	// Subscribe the address before fetching transactions
	parser.Subscribe("0xtestaddress")

	// Fetch transactions for the subscribed address
	transactions := parser.GetTransactions("0xtestaddress")

	// Debugging log to see the number of transactions fetched
	log.Debug.Printf("Number of transactions fetched: %d", len(transactions))

	// Check that we fetched 0 transactions
	assert.Len(t, transactions, 0, "Should return 0 transactions")
}

// Test recording transactions and avoiding duplicates
func TestIsRecordedAndRecordTransaction(t *testing.T) {
	log := logger.GetLogger("debug")
	client := &mockRPCClient{}
	storage := storage.NewMemoryStorage()

	parser := NewEthParser(client, storage, log)

	// Test recording a transaction
	parser.recordTransaction("0x1")
	assert.True(t, parser.isRecorded("0x1"), "Transaction should be recorded")

	// Test that a new transaction is not recorded
	assert.False(t, parser.isRecorded("0x2"), "New transaction should not be recorded yet")
}
