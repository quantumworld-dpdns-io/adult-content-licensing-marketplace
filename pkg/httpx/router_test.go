package httpx

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthz(t *testing.T) {
	t.Parallel()
	mux := NewMux()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var got HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Status != "ok" || got.Service != "gateway" {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestReadyz(t *testing.T) {
	t.Parallel()
	mux := NewMux()
	req := httptest.NewRequest(http.MethodGet, "/readyz", nil)
	rr := httptest.NewRecorder()

	mux.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var got HealthResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if got.Status != "ready" || got.Service != "gateway" {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestCreateAndListLicense(t *testing.T) {
	t.Parallel()
	mux := NewMux()

	createBody := []byte(`{"id":"lic_1","creator_id":"u_1","title":"Sample","description":"d","ai_training_prohibited":true,"base_price_cents":1000}`)
	createReq := httptest.NewRequest(http.MethodPost, "/v1/licenses", bytes.NewReader(createBody))
	createRR := httptest.NewRecorder()
	mux.ServeHTTP(createRR, createReq)
	if createRR.Code != http.StatusCreated {
		t.Fatalf("expected created, got %d", createRR.Code)
	}

	listReq := httptest.NewRequest(http.MethodGet, "/v1/licenses", nil)
	listRR := httptest.NewRecorder()
	mux.ServeHTTP(listRR, listReq)
	if listRR.Code != http.StatusOK {
		t.Fatalf("expected ok, got %d", listRR.Code)
	}

	var out map[string][]map[string]any
	if err := json.Unmarshal(listRR.Body.Bytes(), &out); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if len(out["items"]) != 1 {
		t.Fatalf("expected 1 item, got %d", len(out["items"]))
	}
}
