package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const HeaderRequestID = "X-Request-ID"

func WithRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := r.Header.Get(HeaderRequestID)
		if rid == "" {
			rid = genRequestID()
		}
		w.Header().Set(HeaderRequestID, rid)
		r.Header.Set(HeaderRequestID, rid)
		next.ServeHTTP(w, r)
	})
}

func genRequestID() string {
	b := make([]byte, 12)
	_, err := rand.Read(b)
	if err != nil {
		return "rid-fallback"
	}
	return hex.EncodeToString(b)
}
