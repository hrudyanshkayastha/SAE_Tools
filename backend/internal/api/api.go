package api

import (
	"encoding/json"
	"net/http"
	"sae-core/internal/storage"
)

type Server struct {
	store *storage.Storage
}

func NewServer(store *storage.Storage) *Server {
	return &Server{store: store}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/ready", s.handleReady)
	mux.HandleFunc("/events", s.handleEvents)
	mux.HandleFunc("/incidents", s.handleIncidents)
	mux.HandleFunc("/investigations", s.handleInvestigations)
	mux.HandleFunc("/decisions", s.handleDecisions)

	corsMux := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	}

	return http.ListenAndServe(addr, corsMux(mux))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func (s *Server) handleReady(w http.ResponseWriter, r *http.Request) {
	if s.store.PG.Ping() != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte(`{"status":"postgres unavailable"}`))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ready"}`))
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	rows, err := s.store.PG.Query("SELECT raw FROM sae_telemetry ORDER BY timestamp DESC LIMIT 50")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []json.RawMessage
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err == nil {
			events = append(events, json.RawMessage(raw))
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}

func (s *Server) handleIncidents(w http.ResponseWriter, r *http.Request) {
	rows, err := s.store.PG.Query("SELECT json_build_object('correlation_id', correlation_id, 'target', target, 'events_count', events_count, 'severity', severity, 'status', status) FROM correlations")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []json.RawMessage
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err == nil {
			items = append(items, json.RawMessage(raw))
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

func (s *Server) handleInvestigations(w http.ResponseWriter, r *http.Request) {
	// Investigations correspond to our correlation engine logs which we persist as decisions.
	s.handleDecisions(w, r)
}

func (s *Server) handleDecisions(w http.ResponseWriter, r *http.Request) {
	rows, err := s.store.PG.Query("SELECT json_build_object('correlation_id', correlation_id, 'risk_score', risk_score, 'validation', validation, 'action', action) FROM decisions")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var items []json.RawMessage
	for rows.Next() {
		var raw string
		if err := rows.Scan(&raw); err == nil {
			items = append(items, json.RawMessage(raw))
		}
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(items)
}

