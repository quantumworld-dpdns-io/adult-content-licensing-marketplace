package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/audit"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/auth"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/license"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/policy"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/pkg/payment"
)

type LicenseHandler struct {
	Repo  license.Repository
	Audit audit.Store
}

func (h LicenseHandler) List(w http.ResponseWriter, r *http.Request) {
	principal, err := auth.PrincipalFromRequest(r)
	if err != nil || !auth.CanListLicenses(principal) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	items, err := h.Repo.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "list failed"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (h LicenseHandler) GetByID(w http.ResponseWriter, r *http.Request, id string) {
	principal, err := auth.PrincipalFromRequest(r)
	if err != nil || !auth.CanListLicenses(principal) {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	item, err := h.Repo.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, license.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		if errors.Is(err, license.ErrConflict) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "version conflict"})
			return
		}
		if errors.Is(err, license.ErrBadVersion) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "version is required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "get failed"})
		return
	}
	if h.Audit != nil {
		_ = h.Audit.Append(r.Context(), audit.Event{Type: "license.read", ActorSub: principal.Sub, ActorRole: principal.Role, EntityID: id})
	}
	writeJSON(w, http.StatusOK, item)
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
	created, err := h.Repo.Create(r.Context(), in)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "create failed"})
		return
	}
	if h.Audit != nil {
		_ = h.Audit.Append(r.Context(), audit.Event{Type: "license.created", ActorSub: principal.Sub, ActorRole: principal.Role, EntityID: created.ID, Meta: map[string]any{"currency": created.Currency, "tier": principal.Tier}})
	}
	writeJSON(w, http.StatusCreated, created)
}

func (h LicenseHandler) UpdateByID(w http.ResponseWriter, r *http.Request, id string) {
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
	if in.CreatorID == "" || in.Title == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "creator_id and title are required"})
		return
	}
	if in.Version <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "version is required"})
		return
	}
	if in.Currency == "" {
		in.Currency = "USDC"
	}
	if err := payment.EnsureCryptoCurrency(in.Currency); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	updated, err := h.Repo.Update(r.Context(), id, in)
	if err != nil {
		if errors.Is(err, license.ErrNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}
		if errors.Is(err, license.ErrConflict) {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "version conflict"})
			return
		}
		if errors.Is(err, license.ErrBadVersion) {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "version is required"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "update failed"})
		return
	}
	if h.Audit != nil {
		_ = h.Audit.Append(r.Context(), audit.Event{Type: "license.updated", ActorSub: principal.Sub, ActorRole: principal.Role, EntityID: id, Meta: map[string]any{"currency": updated.Currency, "tier": principal.Tier}})
	}
	writeJSON(w, http.StatusOK, updated)
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
