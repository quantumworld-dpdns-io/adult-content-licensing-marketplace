package api

import (
	"net/http"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/auth"
	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/metrics"
)

func LatencyStats(store *metrics.LatencyStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		principal, err := auth.PrincipalFromRequest(r)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
			return
		}
		if principal.Role != "admin" && principal.Role != "auditor" {
			writeJSON(w, http.StatusForbidden, map[string]string{"error": "forbidden"})
			return
		}
		items := []metrics.EndpointStats{}
		if store != nil {
			items = store.Snapshot()
		}
		writeJSON(w, http.StatusOK, map[string]any{"items": items})
	}
}
