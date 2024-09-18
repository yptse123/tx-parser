// internal/interfaces/interfaces.go
package interfaces

type Parser interface {
	GetCurrentBlock() int
	Subscribe(address string) bool
	GetTransactions(address string) []Transaction
}

type Storage interface {
	AddAddress(address string) bool
	GetTransactions(address string) []Transaction
	AddTransaction(address string, tx Transaction)
}

type Transaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Incoming bool   `json:"incoming"`
}
