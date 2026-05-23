package api

import (
	"net/http"
	"os"
	"time"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/auth"
)

func DevToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.NotFound(w, r)
		return
	}

	role := r.URL.Query().Get("role")
	tier := r.URL.Query().Get("tier")
	sub := r.URL.Query().Get("sub")
	if sub == "" {
		sub = "dev-user"
	}

	p, err := auth.PrincipalFromRequest(&http.Request{Header: http.Header{
		"X-Role": []string{role},
		"X-Tier": []string{tier},
		"X-Sub":  []string{sub},
	}})
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid role/tier"})
		return
	}

	secret := os.Getenv("AUTH_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}
	token, err := auth.SignToken(secret, auth.Claims{
		Sub:  p.Sub,
		Role: p.Role,
		Tier: p.Tier,
		Exp:  time.Now().Add(24 * time.Hour).Unix(),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "token generation failed"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"token": token})
}
