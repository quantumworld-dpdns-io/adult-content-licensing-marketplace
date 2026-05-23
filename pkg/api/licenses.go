package api

import (
	"encoding/json"
	"net/http"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"
)

type LicenseHandler struct {
	Store *license.Store
}

func (h LicenseHandler) List(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{"items": h.Store.List()})
}

func (h LicenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in license.License
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}
	if in.ID == "" || in.CreatorID == "" || in.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "id, creator_id, and title are required"})
		return
	}
	if in.Currency == "" {
		in.Currency = "USD"
	}
	created := h.Store.Create(in)
	writeJSON(w, http.StatusCreated, created)
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
