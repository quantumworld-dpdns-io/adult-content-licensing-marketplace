package api

import (
	"net/http"
	"strconv"
	"time"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/audit"
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
		writeJSON(w, http.StatusOK, map[string]any{"items": []any{}, "total": 0, "limit": 100, "offset": 0, "next_cursor": ""})
		return
	}

	q, err := parseAuditQuery(r)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	n := q.Normalized()

	total, err := h.Audit.Count(r.Context(), n)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "audit count failed"})
		return
	}
	items, err := h.Audit.Query(r.Context(), n)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "audit list failed"})
		return
	}
	nextCursor, err := h.Audit.NextCursor(r.Context(), n)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items":       items,
		"total":       total,
		"limit":       n.Limit,
		"offset":      n.Offset,
		"next_cursor": nextCursor,
	})
}

func parseAuditQuery(r *http.Request) (audit.Query, error) {
	vals := r.URL.Query()
	q := audit.Query{Type: vals.Get("type"), EntityID: vals.Get("entity_id"), Cursor: vals.Get("cursor")}
	if s := vals.Get("since"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return audit.Query{}, err
		}
		q.Since = &t
	}
	if s := vals.Get("until"); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return audit.Query{}, err
		}
		q.Until = &t
	}
	if s := vals.Get("limit"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return audit.Query{}, err
		}
		q.Limit = n
	}
	if s := vals.Get("offset"); s != "" {
		n, err := strconv.Atoi(s)
		if err != nil {
			return audit.Query{}, err
		}
		q.Offset = n
	}
	return q, nil
}
