package storage

import (
	"strings"
	"sync"
	"tx-parser/internal/interfaces"
)

// Normalize address function for consistency
func normalizeAddress(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}

type MemoryStorage struct {
	mu           sync.RWMutex
	subscribed   map[string]bool
	transactions map[string][]interfaces.Transaction
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		subscribed:   make(map[string]bool),
		transactions: make(map[string][]interfaces.Transaction),
	}
}

func (s *MemoryStorage) AddAddress(address string) bool {
	address = normalizeAddress(address)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.subscribed[address]; exists {
		return false
	}

	s.subscribed[address] = true
	return true
}

func (s *MemoryStorage) GetTransactions(address string) []interfaces.Transaction {
	address = normalizeAddress(address)

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.transactions[address]
}

func (s *MemoryStorage) AddTransaction(address string, tx interfaces.Transaction) {
	address = normalizeAddress(address)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Append the transaction to the address's transaction history
	s.transactions[address] = append(s.transactions[address], tx)
}
