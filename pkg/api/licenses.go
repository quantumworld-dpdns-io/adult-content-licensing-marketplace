package api

import (
	"encoding/json"
	"net/http"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/auth"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/policy"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/pkg/payment"
)

type LicenseHandler struct {
	Store *license.Store
}

func (h LicenseHandler) List(w http.ResponseWriter, r *http.Request) {
	principal, err := auth.PrincipalFromRequest(r)
	if err != nil || !auth.CanListLicenses(principal) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": h.Store.List()})
}

func (h LicenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	principal, err := auth.PrincipalFromRequest(r)
	if err != nil || !auth.CanCreateLicense(principal) {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}

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
		in.Currency = "USDC"
	}
	if err := payment.EnsureCryptoCurrency(in.Currency); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	created := h.Store.Create(in)
	writeJSON(w, http.StatusCreated, created)
}

func (h LicenseHandler) CompliancePolicy(w http.ResponseWriter, r *http.Request) {
	principal, err := auth.PrincipalFromRequest(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if principal.Role != "admin" && principal.Role != "auditor" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}
	writeJSON(w, http.StatusOK, policy.Compliance())
}

func (h LicenseHandler) ThreatModelPolicy(w http.ResponseWriter, r *http.Request) {
	principal, err := auth.PrincipalFromRequest(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if principal.Role != "admin" && principal.Role != "auditor" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}
	writeJSON(w, http.StatusOK, policy.ThreatModel())
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}
