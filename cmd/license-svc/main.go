package main

import (
	"log"
	"net/http"
	"os"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/audit"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/config"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/pkg/api"
)

func main() {
	cfg := config.Load()
	_ = os.Setenv("AUTH_SECRET", cfg.AuthSecret)

	repo, err := license.NewRepository(cfg.LicenseStoreBackend, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("repository init failed: %v", err)
	}

	h := api.LicenseHandler{Repo: repo, Audit: audit.NewStream()}
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/auth/dev-token", api.DevToken)
	mux.HandleFunc("/v1/licenses", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			h.List(w, r)
		case http.MethodPost:
			h.Create(w, r)
		default:
			http.NotFound(w, r)
		}
	})
	mux.HandleFunc("/v1/audit/events", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		h.ListAuditEvents(w, r)
	})

	log.Printf("license-svc listening on :8081 (license backend=%s)", cfg.LicenseStoreBackend)
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}
}
