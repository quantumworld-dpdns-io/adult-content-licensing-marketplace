package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/quantumworld-dpdns-io/adult-content-licensing-marketplace/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (w *statusRecorder) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func WithObservability(store *metrics.LatencyStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		d := time.Since(start).Milliseconds()
		if store != nil {
			store.Observe(r.Method+" "+r.URL.Path, d)
		}
		log.Printf("request method=%s path=%s status=%d duration_ms=%d request_id=%s", r.Method, r.URL.Path, rec.status, d, r.Header.Get(HeaderRequestID))
	})
}
