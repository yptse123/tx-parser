package storage

import (
	"testing"
	"tx-parser/internal/interfaces"

	"github.com/stretchr/testify/assert"
)

func TestAddAddress(t *testing.T) {
	storage := NewMemoryStorage()

	// Test adding a new address
	address := "0xTestAddress"
	success := storage.AddAddress(address)
	assert.True(t, success, "Address should be added successfully")

	// Test adding the same address again
	success = storage.AddAddress(address)
	assert.False(t, success, "Adding the same address again should fail")

	// Test adding the same address with different case and whitespace
	addressWithDifferentFormat := "  0xTESTADDRESS "
	success = storage.AddAddress(addressWithDifferentFormat)
	assert.False(t, success, "Adding the same address with different format should fail")
}

func TestGetTransactions(t *testing.T) {
	storage := NewMemoryStorage()

	// Add a few transactions for an address
	address := "0xTestAddress"
	tx1 := interfaces.Transaction{Hash: "0x1", From: "0xFrom1", To: "0xTestAddress", Value: "100"}
	tx2 := interfaces.Transaction{Hash: "0x2", From: "0xTestAddress", To: "0xTo1", Value: "200"}

	// Initially, there should be no transactions
	transactions := storage.GetTransactions(address)
	assert.Len(t, transactions, 0, "There should be no transactions initially")

	// Add transactions to the address
	storage.AddTransaction(address, tx1)
	storage.AddTransaction(address, tx2)

	// Get the transactions and check their content
	transactions = storage.GetTransactions(address)
	assert.Len(t, transactions, 2, "There should be 2 transactions for the address")

	// Check the content of the first transaction
	assert.Equal(t, "0x1", transactions[0].Hash, "The first transaction's hash should match")
	assert.Equal(t, "0xTestAddress", transactions[0].To, "The first transaction's 'To' field should match")
	assert.Equal(t, "100", transactions[0].Value, "The first transaction's value should match")

	// Check the content of the second transaction
	assert.Equal(t, "0x2", transactions[1].Hash, "The second transaction's hash should match")
	assert.Equal(t, "0xTestAddress", transactions[1].From, "The second transaction's 'From' field should match")
	assert.Equal(t, "200", transactions[1].Value, "The second transaction's value should match")
}

func TestAddTransaction_NormalizedAddress(t *testing.T) {
	storage := NewMemoryStorage()

	// Add a transaction to an address with mixed case and extra whitespace
	address := "  0xTestAddress  "
	tx := interfaces.Transaction{Hash: "0x1", From: "0xFrom1", To: "0xTestAddress", Value: "100"}

	// Add the transaction
	storage.AddTransaction(address, tx)

	// Retrieve the transactions using the normalized address
	transactions := storage.GetTransactions("0xtestaddress")
	assert.Len(t, transactions, 1, "There should be 1 transaction after normalization")

	// Check the content of the transaction
	assert.Equal(t, "0x1", transactions[0].Hash, "The transaction's hash should match")
	assert.Equal(t, "0xTestAddress", transactions[0].To, "The transaction's 'To' field should match")
	assert.Equal(t, "100", transactions[0].Value, "The transaction's value should match")
}
