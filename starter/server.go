package main

import (
	"encoding/json"
	"net/http"
)

// Server wraps the HTTP handlers with dependencies
type Server struct {
	store PurchaseStore
}

// NewServer creates a new HTTP server with routes
func NewServer(store PurchaseStore) http.Handler {
	s := &Server{store: store}
	
	mux := http.NewServeMux()
	
	// TODO: Add middleware (logging, request ID, etc.)
	mux.HandleFunc("POST /ingest", s.handleIngest)
	mux.HandleFunc("GET /purchases", s.handleListPurchases)
	
	
	return mux
}

// IngestResponse represents the response from file ingestion
type IngestResponse struct {
	Created int `json:"created"`
	Updated int `json:"updated"`
	Total   int `json:"total"`
}

// ListPurchasesResponse represents the response from listing purchases
type ListPurchasesResponse struct {
	Purchases   []Purchase `json:"purchases"`
	NextAfterID int64      `json:"next_after_id,omitempty"`
}

// handleIngest processes file uploads (NDJSON only)
func (s *Server) handleIngest(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse multipart form file upload
	// Stream parse the NDJSON file without loading it all into memory
	// Use StreamNDJSON parser
	// Track created vs updated counts
	// Return IngestResponse as JSON
	
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing or invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	// TODO: Call StreamNDJSON parser
	// TODO: Count operations and return response
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(IngestResponse{
		Created: 0,
		Updated: 0,
		Total:   0,
	})
}

// handleListPurchases implements keyset pagination for purchases
func (s *Server) handleListPurchases(w http.ResponseWriter, r *http.Request) {
	// TODO: Parse query parameters
	// - after_id: int64 (default 0)
	// - limit: int (default 20, max 100)
	// Call store.ListAfterID with proper context timeout
	// Return ListPurchasesResponse as JSON
	
	_ = int64(0) // TODO: Parse after_id parameter
	// afterID := int64(0)
	// if after := r.URL.Query().Get("after_id"); after != "" {
	//	var err error
	//	afterID, err = strconv.ParseInt(after, 10, 64)
	//	if err != nil {
	//		http.Error(w, "Invalid after_id parameter", http.StatusBadRequest)
	//		return
	//	}
	// }
	
	_ = 20 // TODO: Parse limit parameter
	// limit := 20
	// if l := r.URL.Query().Get("limit"); l != "" {
	//	var err error
	//	limit, err = strconv.Atoi(l)
	//	if err != nil || limit < 1 || limit > 100 {
	//		http.Error(w, "Invalid limit parameter (1-100)", http.StatusBadRequest)
	//		return
	//	}
	// }
	
	// TODO: Call store.ListAfterID with context timeout
	// TODO: Return JSON response
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ListPurchasesResponse{
		Purchases:   []Purchase{},
		NextAfterID: 0,
	})
}


// Helper function to write JSON error responses
func writeJSONError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}