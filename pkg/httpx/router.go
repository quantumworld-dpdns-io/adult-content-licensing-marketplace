package httpx

import (
	"encoding/json"
	"net/http"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/pkg/api"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

func NewMux(repo license.Repository) *http.ServeMux {
	mux := http.NewServeMux()
	licenses := api.LicenseHandler{Repo: repo}
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/readyz", readyzHandler)
	mux.HandleFunc("/v1/auth/dev-token", api.DevToken)
	mux.HandleFunc("/v1/licenses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			licenses.List(w, r)
		case http.MethodPost:
			licenses.Create(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/v1/policy/compliance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		licenses.CompliancePolicy(w, r)
	})
	mux.HandleFunc("/v1/policy/threat-model", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		licenses.ThreatModelPolicy(w, r)
	})
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
