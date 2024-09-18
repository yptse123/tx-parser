package api

import (
	"encoding/json"
	"net/http"

	"tx-parser/internal/interfaces"
	"tx-parser/pkg/logger"
)

type Server struct {
	parser  interfaces.Parser
	log     *logger.Logger
	storage interfaces.Storage
}

func NewServer(p interfaces.Parser, s interfaces.Storage, log *logger.Logger) *Server {
	return &Server{
		parser:  p,
		log:     log,
		storage: s,
	}
}

func (s *Server) Start(address string) error {
	// Attach routes to the default ServeMux
	http.HandleFunc("/subscribe", s.subscribe)
	http.HandleFunc("/transactions/", s.getTransactions) // Route parameter handled manually
	http.HandleFunc("/current-block", s.getCurrentBlock)

	// Start the server and return any error that occurs
	err := http.ListenAndServe(address, nil)
	if err != nil {
		s.log.Error.Printf("Server failed to start: %v", err)
		return err
	}
	return nil
}

func (s *Server) getCurrentBlock(w http.ResponseWriter, r *http.Request) {
	block := s.parser.GetCurrentBlock()
	s.log.Debug.Printf("Fetching current block: %d", block)
	json.NewEncoder(w).Encode(map[string]int{"current_block": block})
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Address string `json:"address"`
	}
	json.NewDecoder(r.Body).Decode(&req)

	// Call the Subscribe method of eth_parser to handle the logic
	if s.parser.Subscribe(req.Address) {
		// If address is newly subscribed, return success
		s.log.Info.Printf("Successfully subscribed to address: %s", req.Address)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	} else {
		// If the address is already subscribed, return conflict
		s.log.Warn.Printf("Address %s is already subscribed", req.Address)
		w.WriteHeader(http.StatusConflict) // 409 Conflict
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "Address already subscribed"})
	}
}

func (s *Server) getTransactions(w http.ResponseWriter, r *http.Request) {
	// Extract the address from the URL
	address := r.URL.Path[len("/transactions/"):]

	// Check if the address is empty
	if address == "" {
		http.Error(w, "Address is required", http.StatusBadRequest)
		return
	}

	// Fetch transactions from storage or the mockParser
	transactions := s.parser.GetTransactions(address)

	// If no transactions found, return a 404
	if transactions == nil {
		http.Error(w, "No transactions found for the given address", http.StatusNotFound)
		return
	}

	// Respond with the transactions
	json.NewEncoder(w).Encode(transactions)
}
