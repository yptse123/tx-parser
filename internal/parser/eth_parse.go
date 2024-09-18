package parser

import (
	"sync"
	"tx-parser/internal/interfaces"
	"tx-parser/internal/rpc"
	"tx-parser/pkg/logger"
	"tx-parser/utils"
)

type EthParser struct {
	currentBlock int
	rpcClient    rpc.Client
	storage      interfaces.Storage
	log          *logger.Logger
	recordedTxns map[string]bool // Tracks recorded transactions (transaction hash as key)
	mu           sync.Mutex      // Protects concurrent access to memory
}

func NewEthParser(client rpc.Client, storage interfaces.Storage, log *logger.Logger) *EthParser {
	// Fetch the current block from the RPC client
	blockNumber, err := client.FetchCurrentBlock()
	if err != nil {
		log.Error.Printf("Error fetching current block during initialization: %v", err)
		blockNumber = 0 // Fallback to 0 in case of error
	} else {
		log.Info.Printf("Fetched current block: %d during initialization", blockNumber)
	}

	return &EthParser{
		currentBlock: blockNumber,
		rpcClient:    client,
		storage:      storage,
		log:          log,
		recordedTxns: make(map[string]bool), // Initialize the recorded transactions map
	}
}

// GetCurrentBlock fetches and updates the current block number
func (p *EthParser) GetCurrentBlock() int {
	blockNumber, err := p.rpcClient.FetchCurrentBlock()
	if err != nil {
		p.log.Error.Printf("Error fetching current block: %v", err)
		return p.currentBlock
	}
	p.currentBlock = blockNumber
	p.log.Info.Printf("Current block updated to %d", p.currentBlock)
	return p.currentBlock
}

// Subscribe adds an address to the storage (if not already subscribed)
func (p *EthParser) Subscribe(address string) bool {
	// Normalize the address before subscribing
	address = utils.NormalizeAddress(address)

	if p.storage.AddAddress(address) {
		p.log.Info.Printf("Address %s successfully subscribed", address)
		return true
	}
	p.log.Warn.Printf("Address %s is already subscribed", address)
	return false
}

// Subscribe adds an address to the storage (if not already subscribed)
func (p *EthParser) GetTransactions(address string) []interfaces.Transaction {
	// Normalize the address
	address = utils.NormalizeAddress(address)

	// Fetch the latest block number
	blockNumber, err := p.rpcClient.FetchCurrentBlock()
	if err != nil || blockNumber == 0 {
		p.log.Error.Println("Error fetching block number")
		return nil
	}

	var newTransactions []interfaces.Transaction

	// Iterate through the blocks and filter transactions for the address
	for i := p.currentBlock; i <= blockNumber; i++ {
		block, err := p.rpcClient.FetchBlockByNumber(i)
		if err != nil {
			p.log.Error.Printf("Error fetching block %d: %v", i, err)
			continue
		}

		// Filter transactions for the address (inbound or outbound)
		for _, tx := range block.Transactions {
			from := utils.NormalizeAddress(tx.From)
			to := utils.NormalizeAddress(tx.To)

			// If the address is involved, determine if it's incoming or outgoing
			if from == address {
				// Outgoing transaction
				newTransactions = append(newTransactions, interfaces.Transaction{
					Hash:     tx.Hash,
					From:     from,
					To:       to,
					Value:    tx.Value,
					Incoming: false, // Outgoing
				})
			} else if to == address {
				// Incoming transaction
				newTransactions = append(newTransactions, interfaces.Transaction{
					Hash:     tx.Hash,
					From:     from,
					To:       to,
					Value:    tx.Value,
					Incoming: true, // Incoming
				})
			}

			// Record the transaction to avoid duplicates
			if !p.isRecorded(tx.Hash) {
				p.recordTransaction(tx.Hash)
				// Store the transaction in the memory storage
				p.storage.AddTransaction(address, tx)
			}
		}
	}

	// Update the current block after processing
	p.currentBlock = blockNumber

	p.log.Info.Printf("Fetched %d new transactions for address: %s", len(newTransactions), address)
	return newTransactions
}

// isRecorded checks if a transaction has already been recorded
func (p *EthParser) isRecorded(txHash string) bool {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.recordedTxns[txHash]
}

// recordTransaction records a new transaction hash in memory
func (p *EthParser) recordTransaction(txHash string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.recordedTxns[txHash] = true
}
