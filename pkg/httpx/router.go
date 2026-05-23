package httpx

import (
	"encoding/json"
	"net/http"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/readyz", readyzHandler)
	return mux
}

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{Status: "ok", Service: "gateway"})
}

func readyzHandler(w http.ResponseWriter, _ *http.Request) {
	respondJSON(w, http.StatusOK, HealthResponse{Status: "ready", Service: "gateway"})
}

func respondJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
