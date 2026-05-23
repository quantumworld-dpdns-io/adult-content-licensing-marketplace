package api

import (
	"net/http"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/auth"
)

func (h LicenseHandler) ListAuditEvents(w http.ResponseWriter, r *http.Request) {
	principal, err := auth.PrincipalFromRequest(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if principal.Role != "admin" && principal.Role != "auditor" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
		return
	}
	if h.Audit == nil {
		writeJSON(w, http.StatusOK, map[string]any{"items": []any{}})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": h.Audit.List()})
}
